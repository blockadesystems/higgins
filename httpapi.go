package main

import (
	"encoding/json"
	// "log"
	"net/http"

	bolt "go.etcd.io/bbolt"
)

type RaftAPI struct {
	node *RaftNode
}

func NewRaftAPI(node *RaftNode) *RaftAPI {
	return &RaftAPI{node: node}
}

func (api *RaftAPI) GetValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	// Get the value from the Raft node
	var value string
	err := api.node.raftStorage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("raft"))
		v := b.Get([]byte(key))
		if v == nil {
			return nil
		}

		value = string(v)
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// value, err := api.node.GetValue(key)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	json.NewEncoder(w).Encode(map[string]string{
		"value": value,
	})
}

func (api *RaftAPI) SetValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	// Propose the new value to the Raft node
	api.node.Propose(key, value)

	w.WriteHeader(http.StatusOK)
}

// func main() {
// 	node := NewRaftNode()

// 	api := NewRaftAPI(node)

// 	http.HandleFunc("/get", api.GetValue)
// 	http.HandleFunc("/set", api.SetValue)

// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
