package main

import (
	"context"
	"encoding/json"
	"log"

	bolt "go.etcd.io/bbolt"
	"go.etcd.io/raft/v3"
)

type Proposal struct {
	K string
	V string
}

type RaftNode struct {
	node        raft.Node
	raftStorage *Kvstore
}

func NewRaftNode() *RaftNode {
	// Initialize a new BoltDB database
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize a new Kvstore
	storage := NewKvstore(db)

	// Create a configuration for the Raft node
	c := &raft.Config{
		ID:              1, // This should be a unique ID for each Raft node
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}

	// Create a new Raft node
	node := raft.StartNode(c, nil)

	return &RaftNode{
		node:        node,
		raftStorage: storage,
	}
}

func (rn *RaftNode) Propose(k, v string) {
	// Create a new proposal
	p := &Proposal{K: k, V: v}

	// Marshal the proposal to a byte slice
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("failed to marshal proposal: %v", err)
	}

	// Propose the command to the Raft node
	if err := rn.node.Propose(context.Background(), data); err != nil {
		log.Fatalf("failed to propose command: %v", err)
	}
}
