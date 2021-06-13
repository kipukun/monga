package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const mangaURL = "https://api.mangadex.org/manga/%s/feed?order[chapter]=desc&limit=500&translatedLanguage[]=en"

type state struct {
	c   *http.Client
	srv *http.Server
}

func (s *state) getIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("/manga/$id"))
}

func (s *state) getFeed(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id, ok := v["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var feedresp FeedResponse
	err := get(r.Context(), s.c, "GET", fmt.Sprintf(mangaURL, id), &feedresp, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rfeed, err := jsonToRss(r.Context(), s.c, id, &feedresp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rss, err := rfeed.ToRss()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, rss)
}

func main() {
	s := new(state)
	s.c = &http.Client{
		Timeout: 10 * time.Second,
	}
	r := mux.NewRouter()
	r.HandleFunc("/", s.getIndex)
	r.HandleFunc("/manga/{id}", s.getFeed)
	s.srv = &http.Server{
		Addr:         ":8000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	closed := make(chan struct{})

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		fmt.Println("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := s.srv.Shutdown(ctx)
		if err != nil {
			fmt.Println("shutdown error:", err)
		}
		cancel()
		close(closed)
	}()

	fmt.Println("listening on 8000")
	err := s.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalln(err)
	}

	<-closed

}
