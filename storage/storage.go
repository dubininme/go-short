package storage

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"sync"
)

type record struct {
	Key string
	URL string
}

type UrlStorage struct {
	urls   map[string]string
	mu     sync.RWMutex
	file   *os.File
	keyGen Generator
}

func NewUrlStorage(keyGen Generator, filename string) *UrlStorage {
	s := &UrlStorage{urls: map[string]string{}, keyGen: keyGen}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening URLStorage:", err)
	}

	s.file = f
	if err := s.load(); err != nil {
		log.Fatal("Error loading URLStorage:", err)
	}
	return s
}

func (s *UrlStorage) Get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.urls[key]
}

func (s *UrlStorage) Set(key, url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, present := s.urls[key]; present {
		return false
	}

	s.urls[key] = url
	return true
}

func (s *UrlStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.urls)
}

func (s *UrlStorage) Put(url string) string {
	for {
		key := s.keyGen.Generate(s.Count())
		if ok := s.Set(key, url); ok {
			if err := s.save(key, url); err != nil {
				log.Println("Error saving to URLStorage:", err)
			}
			return key
		}
	}

	panic("shouldn't get here")
}

func (s *UrlStorage) save(key, url string) error {
	encoder := gob.NewEncoder(s.file)
	return encoder.Encode(record{key, url})
}

func (s *UrlStorage) load() error {
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}

	decoder := gob.NewDecoder(s.file)
	var err error
	for err == nil {
		var rec record
		if err = decoder.Decode(&rec); err == nil {
			s.Set(rec.Key, rec.URL)
		}
	}

	if err == io.EOF {
		return nil
	}

	return err
}
