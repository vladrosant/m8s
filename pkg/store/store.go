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

func (s *Store) CreatePod(pod api.Pod) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.state.Pods {
		if p.Name == pod.Name && p.Namespace == pod.Namespace {
			return fmt.Errorf("pod %s/%s already exists", pod.Namespace, pod.Name)
		}
	}

	s.state.Pods = append(s.state.Pods, pod)
	return s.save()
}

func (s *Store) GetPod(namespace, name string) (*api.Pod, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, pod := range s.state.Pods {
		if pod.Name == name && pod.Namespace == namespace {
			podCopy := pod
			return &podCopy, nil
		}
	}

	return nil, fmt.Errorf("pod %s/%s not found!", namespace, name)
}

func (s *Store) ListPods() ([]api.Pod, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pods := make([]api.Pod, len(s.state.Pods))
	copy(pods, s.state.Pods)
	return pods, nil
}

func (s *Store) UpdatePod(pod api.Pod) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.state.Pods {
		if p.Name == pod.Name && p.Namespace == pod.Namespace {
			s.state.Pods[i] = pod
			return s.save()
		}
	}

	return fmt.Errorf("pod %s/%s not found!", pod.Namespace, pod.Name)
}

func (s *Store) DeletePod(namespace, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, pod := range s.state.Pods {
		if pod.Name == name && pod.Namespace == namespace {
			s.state.Pods = append(s.state.Pods[:i], s.state.Pods[i+1:]...)
			return s.save()
		}
	}

	return fmt.Errorf("pod %s/%s not found!", namespace, name)

	// add node methods later on (CreateNote, GetNode, ListNodes...)
}
