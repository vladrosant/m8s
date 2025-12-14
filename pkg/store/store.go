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

func (s *Store) CreatePod(node api.Pod) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.state.Pods {
		if p.Name == node.Name && p.Namespace == node.Namespace {
			return fmt.Errorf("node %s/%s already exists", node.Namespace, node.Name)
		}
	}

	s.state.Pods = append(s.state.Pods, node)
	return s.save()
}

func (s *Store) GetPod(namespace, name string) (*api.Pod, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, node := range s.state.Pods {
		if node.Name == name && node.Namespace == namespace {
			nodeCopy := node
			return &nodeCopy, nil
		}
	}

	return nil, fmt.Errorf("node %s/%s not found!", namespace, name)
}

func (s *Store) ListPods() ([]api.Pod, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodes := make([]api.Pod, len(s.state.Pods))
	copy(nodes, s.state.Pods)
	return nodes, nil
}

func (s *Store) UpdatePod(node api.Pod) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.state.Pods {
		if p.Name == node.Name && p.Namespace == node.Namespace {
			s.state.Pods[i] = node
			return s.save()
		}
	}

	return fmt.Errorf("node %s/%s not found!", node.Namespace, node.Name)
}

func (s *Store) DeletePod(namespace, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, node := range s.state.Pods {
		if node.Name == name && node.Namespace == namespace {
			s.state.Pods = append(s.state.Pods[:i], s.state.Pods[i+1:]...)
			return s.save()
		}
	}

	return fmt.Errorf("node %s/%s not found!", namespace, name)
}

// add node methods later on (CreateNode, GetNode, ListNodes...)

func (s *Store) CreateNode(node api.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, n := range s.state.Nodes {
		if n.Name == node.Name {
			return fmt.Errorf("node %s already exists", node.Name)
		}
	}

	s.state.Nodes = append(s.state.Nodes, node)
	return s.save()
}

func (s *Store) GetNode(name string) (*api.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, node := range s.state.Nodes {
		if node.Name == name {
			nodeCopy := node
			return &nodeCopy, nil
		}
	}

	return nil, fmt.Errorf("node %s not found!", name)
}

func (s *Store) ListNodes(node api.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, n := range s.state.Nodes {
		if n.Name == node.Name {
			s.state.Nodes[i] = node
			return s.save()
		}
	}

	return fmt.Errorf("node %s not found!", node.Name)
}
