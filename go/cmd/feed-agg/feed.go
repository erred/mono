package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

type FeedUpdater struct {
	id     string
	u      string
	client *http.Client
	ch     chan *feeds.Feed

	etag   string
	etagmu sync.Mutex
}

func RunFeedUpdater(id, u string, interval time.Duration, store Storer) {
	fu := &FeedUpdater{
		id:     id,
		u:      u,
		client: http.DefaultClient,
		ch:     make(chan *feeds.Feed),
	}

	go func() {
		fn := func() {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, interval)
			defer cancel()
			err := fu.get(ctx)
			if err != nil {
				log.Println(id, err)
			}
		}

		fn()
		tick := time.NewTicker(interval)
		for range tick.C {
			fn()
		}
	}()

	go func() {
		ctx := context.TODO()
		for f := range fu.ch {
			if f == nil {
				continue
			}
			err := store.UpdateUpstream(ctx, id, f)
			if err != nil {
				log.Println(id, err)
			}
		}
	}()
}

func (f *FeedUpdater) get(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.u, nil)
	if err != nil {
		return err
	}

	f.etagmu.Lock()
	if f.etag != "" {
		req.Header.Set("if-none-match", f.etag)
	}
	f.etagmu.Unlock()

	res, err := f.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotModified {
			f.ch <- nil
			return nil
		}
		return errors.New(res.Status)
	}
	defer res.Body.Close()

	parser := gofeed.NewParser()
	feed, err := parser.Parse(res.Body)
	if err != nil {
		return err
	}

	its := make([]*feeds.Item, len(feed.Items))
	for i := range feed.Items {
		its[i] = convertItem(feed.Items[i])
		its[i].Source = &feeds.Link{
			Href: feed.Title,
		}
	}

	f.ch <- &feeds.Feed{
		Items: its,
	}
	f.etagmu.Lock()
	f.etag = res.Header.Get("etag")
	f.etagmu.Unlock()

	return nil
}

func convertItem(in *gofeed.Item) *feeds.Item {
	f := &feeds.Item{
		Title:       in.Title,
		Description: in.Description,
		Id:          in.GUID,
		Link: &feeds.Link{
			Href: in.Link,
		},
		Content: in.Content,
	}
	if in.Author != nil {
		f.Author = &feeds.Author{
			Name:  in.Author.Name,
			Email: in.Author.Email,
		}
	}
	if in.PublishedParsed != nil {
		f.Created = *in.PublishedParsed
		f.Updated = *in.PublishedParsed // fallback
	}
	if in.UpdatedParsed != nil {
		f.Updated = *in.UpdatedParsed
		if in.PublishedParsed == nil {
			f.Created = *in.UpdatedParsed // fallback
		}
	}
	return f
}
