package main

import (
	"bytes"
	"encoding/binary"

	bolt "go.etcd.io/bbolt"
	"go.etcd.io/raft/v3/raftpb"
)

type Kvstore struct {
	db *bolt.DB
}

func NewKvstore(db *bolt.DB) *Kvstore {
	return &Kvstore{db: db}
}

func (s *Kvstore) InitialState() (raftpb.HardState, raftpb.ConfState, error) {
	var hs raftpb.HardState
	var cs raftpb.ConfState

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("raft"))
		hs = unmarshalHardState(b.Get([]byte("hardstate")))
		cs = unmarshalConfState(b.Get([]byte("confstate")))

		return nil
	})

	return hs, cs, err
}

func (s *Kvstore) Entries(lo, hi, maxSize uint64) ([]raftpb.Entry, error) {
	var entries []raftpb.Entry

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("raft"))
		c := b.Cursor()

		min := itob(lo)
		max := itob(hi)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			entries = append(entries, unmarshalEntry(v))
		}

		return nil
	})

	return entries, err
}

func (s *Kvstore) Term(i uint64) (uint64, error) {
	var term uint64
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("raft"))
		term = binary.BigEndian.Uint64(b.Get(itob(i)))
		return nil
	})

	return term, err
}

func (s *Kvstore) LastIndex() (uint64, error) {
	var lastIndex uint64
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("raft"))
		c := b.Cursor()

		_, v := c.Last()
		lastIndex = binary.BigEndian.Uint64(v)

		return nil
	})

	return lastIndex, err
}

func (s *Kvstore) FirstIndex() (uint64, error) {
	var firstIndex uint64
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("raft"))
		c := b.Cursor()

		_, v := c.First()
		firstIndex = binary.BigEndian.Uint64(v)

		return nil
	})

	return firstIndex, err
}

func (s *Kvstore) Snapshot() (raftpb.Snapshot, error) {
	// Implement this method
	return raftpb.Snapshot{}, nil
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func unmarshalHardState(data []byte) raftpb.HardState {
	var hs raftpb.HardState
	// Unmarshal data into hs
	return hs
}

func unmarshalConfState(data []byte) raftpb.ConfState {
	var cs raftpb.ConfState
	// Unmarshal data into cs
	return cs
}

func unmarshalEntry(data []byte) raftpb.Entry {
	var e raftpb.Entry
	// Unmarshal data into e
	return e
}
