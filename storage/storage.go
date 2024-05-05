package storage

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"sync"
)

const saveQueueLength = 1000

type record struct {
	Key string
	URL string
}

type UrlStorage struct {
	urls   map[string]string
	mu     sync.RWMutex
	keyGen Generator
	save   chan record
}

func NewUrlStorage(keyGen Generator, filename string) *UrlStorage {
	s := &UrlStorage{
		urls:   map[string]string{},
		keyGen: keyGen,
		save:   make(chan record, saveQueueLength),
	}

	if err := s.load(filename); err != nil {
		log.Fatal("Error loading URLStorage: ", err)
	}

	go s.saveLoop(filename)
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
		if s.Set(key, url) {
			s.save <- record{key, url}
			return key
		}
	}

	panic("shouldn't get here")
}

func (s *UrlStorage) saveLoop(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening URLStorage:", err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	for {
		r := <-s.save
		log.Println("saving url to storage in loop")
		if err := encoder.Encode(&r); err != nil {
			log.Println("Error saving to URLStore: ", err)
		}
	}
}

func (s *UrlStorage) load(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error opening URLStorage: ", err)
		return err
	}
	defer f.Close()
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Println("Error seek URLStorage: ", err)
		return err
	}

	decoder := gob.NewDecoder(f)
	var rec record
	for {
		err = decoder.Decode(&rec)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Println("Error decoding rec: ", rec)
			return err
		}

		s.Set(rec.Key, rec.URL)
	}
}
