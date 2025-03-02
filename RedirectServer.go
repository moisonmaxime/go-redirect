package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type RedirectServer struct {
	port       int
	host       string
	storage    *URLStorage
	httpServer *http.Server
}

func newRedirectServer(host string, port int, filename string) (*RedirectServer, error) {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second, // Timeout for reading the request
		WriteTimeout:   10 * time.Second, // Timeout for writing the response
		MaxHeaderBytes: 1 << 20,          // Limit on the size of headers
	}

	urlStorage, err := newURLStorage(filename)
	if err != nil {
		return nil, err
	}

	s := RedirectServer{
		host:       host,
		port:       port,
		storage:    urlStorage,
		httpServer: httpServer,
	}

	mux.HandleFunc("/", s.routeRequests)

	return &s, nil
}

func (server *RedirectServer) start() {
	log.Printf("Server started on :%v", server.port)
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func (server *RedirectServer) stop(shutdownCtx context.Context) {
	if err := server.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	if err := server.storage.SaveToFile(); err != nil {
		log.Fatalf("URLs Save Failed: %v", err)
		log.Fatalf("Current state of URLs (for recovery): %v", server.storage.urls)
	}
}

func (server *RedirectServer) routeRequests(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		server.handleRedirect(w, req)
	case http.MethodPost:
		server.handleCreateShortLink(w, req)
	case http.MethodDelete:
		server.handleDeleteShortLink(w, req)
	}
}

func (server *RedirectServer) handleRedirect(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")

	log.Println(path)

	if len(path) > 0 {
		short := path
		log.Println(short)

		match := server.storage.get(short)
		log.Println(match)

		if match != "" {
			w.Header().Set("Location", match)
			w.WriteHeader(http.StatusPermanentRedirect)
			w.Write([]byte(""))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("\n404: Not Found!\n"))
}

func (server *RedirectServer) handleCreateShortLink(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("url")
	path := strings.TrimPrefix(req.URL.Path, "/")
	shortlink := strings.Split(path, "/")[0]

	if shortlink == "" {
		shortlink = generateRandomString(8)
	}

	server.storage.store(shortlink, url)

	resultMessage := fmt.Sprintf("\n%v was added under shortlink: http://%v:%v/%v\n", url, server.host, server.port, shortlink)
	w.Write([]byte(resultMessage))
}

func (server *RedirectServer) handleDeleteShortLink(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	shortlink := strings.Split(path, "/")[0]

	if server.storage.get(shortlink) == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Not Found!"))
		return
	}

	server.storage.remove(shortlink)

	resultMessage := fmt.Sprintf("\nRemoved shortlink: %v/n", shortlink)
	w.Write([]byte(resultMessage))
}
