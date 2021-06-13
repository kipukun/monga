package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/feeds"
)

var (
	cache sync.Map
	sem   chan struct{}
)

func init() {
	// mangadex only allows 5 requests per second, limit it here
	sem = make(chan struct{}, 5)
}

const cubariURL = "https://cubari.moe/read/mangadex/%s/%s/1"

func get(ctx context.Context, c *http.Client, method, url string, v, data interface{}) error {

	// block until someone is done
	// and release when we return
	sem <- struct{}{}
	defer func() {
		<-sem
	}()

	var err error
	var req *http.Request
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(b))
		if err != nil {
			return err
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return err
		}
	}
	fr, err := c.Do(req)
	if err != nil {
		return err
	}
	feeds, err := io.ReadAll(fr.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(feeds, &v)
}

func jsonToRss(ctx context.Context, c *http.Client, id string, fr *FeedResponse) (*feeds.Feed, error) {
	var items []*feeds.Item
	var title interface{}
	var ok bool
	title, ok = cache.Load(id)
	if !ok {
		var mresp MangaResponse
		err := get(ctx, c, "GET", "https://api.mangadex.org/manga/"+id, &mresp, nil)
		if err != nil {
			return nil, err
		}
		cache.Store(id, mresp.Data.Attributes.Title.En)
		title = mresp.Data.Attributes.Title.En
	}
	for _, res := range fr.Results {
		time, err := time.Parse(time.RFC3339, strings.Replace(res.Data.Attributes.Createdat, "Z00:00", "Z", 1))
		if err != nil {
			return nil, err
		}
		var b strings.Builder
		b.WriteString(title.(string))
		b.WriteString(" - ")
		if res.Data.Attributes.Title != "" {
			b.WriteString(res.Data.Attributes.Title)
			b.WriteString(" - ")
		}

		if res.Data.Attributes.Volume != "" {
			b.WriteString("Vol. ")
			b.WriteString(res.Data.Attributes.Volume)
			b.WriteString(", ")
		}
		b.WriteString("Ch. ")
		b.WriteString(res.Data.Attributes.Chapter)
		i := &feeds.Item{
			Title:   b.String(),
			Link:    &feeds.Link{Href: fmt.Sprintf(cubariURL, id, res.Data.Attributes.Chapter)},
			Created: time,
		}
		items = append(items, i)
	}
	feed := &feeds.Feed{
		Title:   "Mangadex - " + title.(string),
		Link:    &feeds.Link{Href: "https://mangadex.org/title/" + id},
		Created: time.Now(),
		Items:   items,
	}
	return feed, nil
}
