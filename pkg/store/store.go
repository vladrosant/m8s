package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/vladrosant/m8s/pkg/api"
)

type State struct {
	Pods  []api.Pod  `json:"pods"`
	Nodes []api.Node `json:"nodes"`
}

type Store struct {
	filepath string
	mu       sync.RWMutex
	state    State
}

func NewStore(filepath string) (*Store, error) {
	s := &Store{
		filepath: filepath,
		state: State{
			Pods:  []api.Pod{},
			Nodes: []api.Node{},
		},
	}

	if err := s.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load state: %w", err)
		}

		if err := s.save(); err != nil {
			return nil, fmt.Errorf("failed to create state file: %w", err)
		}
	}

	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return json.Unmarshal(data, &s.state)
}

func (s *Store) save() error {
	s.mu.RLock()
	data, err := json.MarshalIndent(s.state, "", "	")
	s.mu.RUnlock()

	if err != nil {
		return err
	}

	return os.WriteFile(s.filepath, data, 0644)
}
