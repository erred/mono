package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"

	_ "modernc.org/sqlite"
)

type Storer interface {
	UpdateUpstream(ctx context.Context, id string, f *feeds.Feed) error
	GetFeed(ctx context.Context, id string) (*feeds.Feed, error)
}

type SQLite struct {
	db *sql.DB
}

func NewSQLite(ctx context.Context, dataDir string, fcs map[string]FeedConfig) (*SQLite, error) {
	db, err := sql.Open("sqlite", filepath.Join(dataDir, "feeds.sqlite"))
	if err != nil {
		return nil, err
	}

	s := &SQLite{
		db: db,
	}

	_, err = s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS upstream (
        id TEXT PRIMARY KEY,
        feed_json TEXT NOT NULL
)
`)
	if err != nil {
		return nil, err
	}
	_, err = s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS feed (
        id TEXT PRIMARY KEY,
        refresh_seconds INTEGER NOT NULL,
        updated TEXT NOT NULL,
        feed_json TEXT NOT NULL
)
`)
	if err != nil {
		return nil, err
	}
	_, err = s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS feed_upstream (
        feed_id TEXT,
        upsteam_id TEXT,
        FOREIGN KEY (feed_id) REFERENCES feed(id),
        FOREIGN KEY (upsteam_id) REFERENCES upstream(id),
        UNIQUE(feed_id, upsteam_id)
)
        `)
	if err != nil {
		return nil, err
	}
	for feedID, feedConfig := range fcs {
		_, err := s.db.ExecContext(ctx, `
REPLACE INTO feed (id, refresh_seconds, updated, feed_json)
VALUES (?, ?, ?, ?)
                `, feedID, int64(feedConfig.refresh.Seconds()), time.Unix(0, 0).Format(time.RFC3339), "")
		if err != nil {
			return nil, err
		}
		for upstreamID := range feedConfig.URLs {
			_, err := s.db.ExecContext(ctx, `
REPLACE INTO upstream (id, feed_json)
VALUES (?, ?)
                `, upstreamID, "")
			if err != nil {
				return nil, err
			}
			_, err = s.db.ExecContext(ctx, `
REPLACE INTO feed_upstream (feed_id, upsteam_id)
VALUES (?, ?)
                        `, feedID, upstreamID)
			if err != nil {
				return nil, err
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
	row := s.db.QueryRowContext(ctx, `
SELECT feed_json, updated, refresh_seconds
FROM feed
WHERE id = ?
        `, id)
	var rawFeed, rawUpdated string
	var refreshSeconds float64
	feed := &feeds.Feed{}

	err := row.Scan(&rawFeed, &rawUpdated, &refreshSeconds)
	if err != nil {
		return nil, err
	}
	updated, err := time.Parse(time.RFC3339, rawUpdated)
	if err != nil {
		return nil, err
	}

	if time.Since(updated).Seconds() > refreshSeconds {
		feed, err = refresh(ctx, s.db, id)
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

func refresh(ctx context.Context, db *sql.DB, id string) (*feeds.Feed, error) {
	rows, err := db.QueryContext(ctx, `
SELECT feed_json
FROM upstream
WHERE id IN (
        SELECT upsteam_id
        FROM feed_upstream
        WHERE feed_id = ?
)
                `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []*feeds.Item
	for rows.Next() {
		var raw string
		err := rows.Scan(&raw)
		if err != nil {
			return nil, err
		}
		var feed feeds.Feed
		err = json.Unmarshal([]byte(raw), &feed)
		if err != nil {
			return nil, err
		}
		all = append(all, feed.Items...)
	}
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

	_, err = db.ExecContext(ctx, `
UPDATE feed
SET
        feed_json = ?,
        updated = ?
WHERE id = ?
        `, string(b), f.Updated.Format(time.RFC3339), id)
	if err != nil {
		return nil, err
	}

	return f, nil
}
