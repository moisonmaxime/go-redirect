package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type URLStorage struct {
	filename string
	urls     map[string]string
	lock     sync.Mutex
	fileLock sync.Mutex
}

func newURLStorage(filename string) (*URLStorage, error) {
	s := URLStorage{filename: filename}

	if err := s.LoadFromFile(); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *URLStorage) SaveToFile() error {
	jsonData, err := json.MarshalIndent(s.urls, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to format urls into json: %v", err)
	}
	err = os.WriteFile(s.filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save urls to file (%v): %v", s.filename, err)
	}

	return nil
}

func (s *URLStorage) LoadFromFile() error {
	b, err := os.ReadFile(s.filename)
	if err != nil {
		return fmt.Errorf("failed to load urls from file (%v): %v", s.filename, err)
	}
	json.Unmarshal(b, &s.urls)

	return nil
}

func (s *URLStorage) store(short string, full string) {
	s.lock.Lock()
	s.urls[short] = full
	s.lock.Unlock()

	// Save to file async
	go func() {
		log.Println("Saving to file async")
		s.fileLock.Lock()
		s.SaveToFile()
		s.fileLock.Unlock()
		log.Println("Done aving to file async")
	}()
}

func (s *URLStorage) get(short string) string {
	return s.urls[short]
}

func (s *URLStorage) remove(short string) {
	s.lock.Lock()
	delete(s.urls, short)
	s.lock.Unlock()

	// Save to file async
	go func() {
		log.Println("Saving to file async")
		s.fileLock.Lock()
		s.SaveToFile()
		s.fileLock.Unlock()
		log.Println("Done aving to file async")
	}()
}
