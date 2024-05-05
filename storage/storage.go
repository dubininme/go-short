package storage

import (
	"encoding/json"
	"errors"
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
	if filename == "" {
		log.Fatal("Error loading URLStorage: filename is empty")
	}

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

func (s *UrlStorage) Get(key, url *string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if cachedUrl, ok := s.urls[*key]; ok {
		*url = cachedUrl
		return nil
	}

	return errors.New("key not found")
}

func (s *UrlStorage) Set(key, url *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, present := s.urls[*key]; present {
		return errors.New("key already exists")
	}

	s.urls[*key] = *url
	return nil
}

func (s *UrlStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.urls)
}

func (s *UrlStorage) Put(url, key *string) error {
	for {
		// TODO refactor if needs rewrite of input key
		*key = s.keyGen.Generate(s.Count())
		if err := s.Set(key, url); err == nil {
			break
		}
	}

	if s.save != nil {
		s.save <- record{*key, *url}
	}

	return nil
}

func (s *UrlStorage) saveLoop(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening URLStorage:", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
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

	decoder := json.NewDecoder(f)
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

		s.Set(&rec.Key, &rec.URL)
	}
}
