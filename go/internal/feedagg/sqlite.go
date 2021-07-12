package feedagg

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	_ "modernc.org/sqlite"
)

type Storer interface {
	UpdateUpstream(ctx context.Context, id string, f *feeds.Feed) error
	GetFeed(ctx context.Context, id string) (*feeds.Feed, error)
}

type SQLite struct {
	db *sql.DB
	t  trace.Tracer
}

func NewSQLite(ctx context.Context, dataDir string, fcs map[string]FeedConfig) (*SQLite, error) {
	tracer := otel.Tracer("sqlite")
	ctx, span := tracer.Start(ctx, "init-sqlite")
	defer span.End()

	db, err := sql.Open("sqlite", filepath.Join(dataDir, "feeds.sqlite"))
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	s := &SQLite{
		db: db,
		t:  tracer,
	}

	lctx, span := tracer.Start(ctx, "create-table", trace.WithAttributes(attribute.String("table", "upstream")))
	_, err = s.db.ExecContext(lctx, `
CREATE TABLE IF NOT EXISTS upstream (
        id TEXT PRIMARY KEY,
        feed_json TEXT NOT NULL
)
`)
	span.End()
	if err != nil {
		return nil, fmt.Errorf("create table upstream: %w", err)
	}

	lctx, span = tracer.Start(ctx, "create-table", trace.WithAttributes(attribute.String("table", "feed")))
	_, err = s.db.ExecContext(lctx, `
CREATE TABLE IF NOT EXISTS feed (
        id TEXT PRIMARY KEY,
        refresh_seconds INTEGER NOT NULL,
        updated TEXT NOT NULL,
        feed_json TEXT NOT NULL
)
`)
	span.End()
	if err != nil {
		return nil, fmt.Errorf("create table feed: %w", err)
	}

	lctx, span = tracer.Start(ctx, "create-table", trace.WithAttributes(attribute.String("table", "feed_upstream")))
	_, err = s.db.ExecContext(lctx, `
CREATE TABLE IF NOT EXISTS feed_upstream (
        feed_id TEXT,
        upsteam_id TEXT,
        FOREIGN KEY (feed_id) REFERENCES feed(id),
        FOREIGN KEY (upsteam_id) REFERENCES upstream(id),
        UNIQUE(feed_id, upsteam_id)
)
        `)
	span.End()
	if err != nil {
		return nil, fmt.Errorf("create table feed_upstream: %w", err)
	}

	for feedID, feedConfig := range fcs {
		lctx, span := tracer.Start(ctx, "create-table", trace.WithAttributes(attribute.String("table", "upstream")))
		_, err := s.db.ExecContext(lctx, `
REPLACE INTO feed (id, refresh_seconds, updated, feed_json)
VALUES (?, ?, ?, ?)
                `, feedID, int64(feedConfig.refresh.Seconds()), time.Unix(0, 0).Format(time.RFC3339), "")
		span.End()
		if err != nil {
			return nil, fmt.Errorf("replace into feed %s: %w", feedID, err)
		}
		for upstreamID := range feedConfig.URLs {
			lctx, span = tracer.Start(ctx, "create-table", trace.WithAttributes(attribute.String("table", "upstream")))
			_, err := s.db.ExecContext(lctx, `
REPLACE INTO upstream (id, feed_json)
VALUES (?, ?)
                `, upstreamID, "")
			span.End()
			if err != nil {
				return nil, fmt.Errorf("replace into upstream %s: %w", feedID, err)
			}

			lctx, span = tracer.Start(ctx, "create-table", trace.WithAttributes(attribute.String("table", "upstream")))
			_, err = s.db.ExecContext(lctx, `
REPLACE INTO feed_upstream (feed_id, upsteam_id)
VALUES (?, ?)
                        `, feedID, upstreamID)
			span.End()
			if err != nil {
				return nil, fmt.Errorf("replace into feed_upstream %s: %w", feedID, err)
			}
		}
	}

	return s, nil
}

func (s *SQLite) UpdateUpstream(ctx context.Context, id string, f *feeds.Feed) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
UPDATE upstream
SET feed_json = ?
WHERE id = ?
        `, string(b), id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLite) GetFeed(ctx context.Context, id string) (*feeds.Feed, error) {
	ctx, span := s.t.Start(ctx, "get-feed", trace.WithAttributes(attribute.String("feed_id", id)))
	defer span.End()

	var rawFeed, rawUpdated string
	var refreshSeconds float64
	feed := &feeds.Feed{}

	ctx, span = s.t.Start(ctx, "query")
	row := s.db.QueryRowContext(ctx, `
SELECT feed_json, updated, refresh_seconds
FROM feed
WHERE id = ?
        `, id)
	err := row.Scan(&rawFeed, &rawUpdated, &refreshSeconds)
	span.End()
	if err != nil {
		return nil, err
	}
	updated, err := time.Parse(time.RFC3339, rawUpdated)
	if err != nil {
		return nil, err
	}

	if time.Since(updated).Seconds() > refreshSeconds {
		feed, err = s.refresh(ctx, s.db, id)
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal([]byte(rawFeed), feed)
		if err != nil {
			return nil, err
		}
	}
	return feed, nil
}

func (s *SQLite) refresh(ctx context.Context, db *sql.DB, id string) (*feeds.Feed, error) {
	ctx, span := s.t.Start(ctx, "refresh", trace.WithAttributes(attribute.String("feed_id", id)))
	defer span.End()

	ctx, span = s.t.Start(ctx, "query")
	rows, err := db.QueryContext(ctx, `
SELECT feed_json
FROM upstream
WHERE id IN (
        SELECT upsteam_id
        FROM feed_upstream
        WHERE feed_id = ?
)
                `, id)
	span.End()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ctx, span = s.t.Start(ctx, "iterate")
	var all []*feeds.Item
	for rows.Next() {
		var raw string
		err := rows.Scan(&raw)
		if err != nil {
			span.End()
			return nil, err
		}
		var feed feeds.Feed
		err = json.Unmarshal([]byte(raw), &feed)
		if err != nil {
			span.End()
			return nil, err
		}
		all = append(all, feed.Items...)
	}
	span.End()
	f := &feeds.Feed{
		Title: id,
		Link: &feeds.Link{
			Href: "/" + id,
			Rel:  "",
		},
		Description: "aggregated feed for " + id,
		Copyright:   "",
		Updated:     time.Now(),
		Items:       all,
	}
	f.Sort(func(a, b *feeds.Item) bool {
		return a.Updated.After(b.Updated)
	})

	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	ctx, span = s.t.Start(ctx, "update")
	_, err = db.ExecContext(ctx, `
UPDATE feed
SET
        feed_json = ?,
        updated = ?
WHERE id = ?
        `, string(b), f.Updated.Format(time.RFC3339), id)
	span.End()
	if err != nil {
		return nil, err
	}

	return f, nil
}
