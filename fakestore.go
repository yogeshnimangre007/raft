package raft

import "errors"

// Implements the LogStore interface.
// It is an in memory store to test the raft node.
// DONT use for production
type FakeStore struct {
	store []Log
}

func (fs *FakeStore) Empty() bool {
	if s := len(fs.store); s == 0 {
		return true
	}
	return false
}

func (fs *FakeStore) LastIndex() uint64 {
	if fs.Empty() {
		return 0
	}
	return fs.store[len(fs.store)].Index
}

func (fs *FakeStore) LastTerm() uint64 {
	if fs.Empty() {
		return 0
	}
	return fs.store[len(fs.store)].Term
}

func (fs *FakeStore) GetLog(index uint64) (*Log, error) {
	if fs.Empty() {
		return nil, errors.New("store is empty")
	}

	s := len(fs.store)
	for i := 0; i < s; i++ {
		if fs.store[i].Index == index {
			return &fs.store[i], nil
		}
	}

	return nil, errors.New("index not found in store")
}

func (fs *FakeStore) StoreLog(log *Log) error {
	// TODO
	return nil
}

func (fs *FakeStore) DeleteRange(min, max uint64) error {
	// TODO
	return nil
}