package main

import "sync"

type Store struct {
	stockItems []stock
	itemsMutex sync.Mutex
}

func (s *Store) Append(item stock) {
	s.itemsMutex.Lock()
	defer s.itemsMutex.Unlock()
	s.stockItems = append(s.stockItems, item)
}

func (s *Store) GetAll() []stock {
	s.itemsMutex.Lock()
	defer s.itemsMutex.Unlock()
	return s.stockItems
}

func (s *Store) Pop() (stock, bool) {
	s.itemsMutex.Lock()
	defer s.itemsMutex.Unlock()
	if len(s.stockItems) == 0 {
		return stock{}, false
	}
	item := s.stockItems[0]
	s.stockItems = s.stockItems[1:]
	return item, true
}
