package cache

import (
	"container/list"
)

// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LRU struct {
	// whatever fields you want here
	Limit           int
	OrderedElements *list.List
	ElementsMap     map[string]*list.Element
	CurrentStorage  int
	Statistics      *Stats
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	return &LRU{Limit: limit, OrderedElements: list.New(), ElementsMap: make(map[string]*list.Element), Statistics: &Stats{}}
}

// MaxStorage returns the maximum number of bytes this LRU can store
func (lru *LRU) MaxStorage() int {
	return lru.Limit
}

// RemainingStorage returns the number of unused bytes available in this LRU
func (lru *LRU) RemainingStorage() int {
	return lru.Limit - lru.CurrentStorage
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {
	existingElement, exists := lru.ElementsMap[key]
	if !exists {
		lru.Statistics.Misses += 1
		return nil, false
	}

	lru.OrderedElements.MoveToFront(existingElement)
	lru.Statistics.Hits += 1
	return existingElement.Value.(ListElement).Value, true
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	existingElement, exists := lru.ElementsMap[key]
	if !exists {
		return nil, false
	}
	lru.OrderedElements.Remove(existingElement)
	lru.CurrentStorage -= len(key) + len(existingElement.Value.(ListElement).Value)
	delete(lru.ElementsMap, key)
	return nil, false
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value []byte) bool {
	insertSize := len(key) + len(value)
	if insertSize > lru.Limit {
		return false
	}

	// remove key if present
	lru.Remove(key)

	for (lru.RemainingStorage() - insertSize) < 0 {
		lruElement := lru.OrderedElements.Back()
		lru.Remove(lruElement.Value.(ListElement).Key)
	}

	lru.ElementsMap[key] = lru.OrderedElements.PushFront(ListElement{Key: key, Value: value})
	lru.CurrentStorage += insertSize
	return true
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return lru.OrderedElements.Len()
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return lru.Statistics
}
