package store

import (
	"os"
	"testing"
	"time"

	"github.com/vladrosant/m8s/pkg/types"
)

func TestStore(t *testing.T) {
	testFile := "/tmp/m8s-test-state.json"
	defer os.Remove(testFile)

	store, err := NewStore(testFile)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	pod := types.Pod{
		Name:      "test-pod",
		Namespace: "default",
		Image:     "nginx.latest",
		Status:    types.PodStatusPending,
		CreatedAt: time.Now(),
	}

	err = store.CreatePod(pod)
	if err != nil {
		t.Fatalf("failed to create pod: %v", err)
	}

	retrieved, err := store.GetPod("default", "test-pod")
	if err != nil {
		t.Fatalf("failed to get pod: %v", err)
	}

	if retrieved.Name != "test-pod" {
		t.Errorf("expected pod name 'test-pod', got '%s'", retrieved.Name)
	}

	pods, err := store.ListPods()
	if err != nil {
		t.Fatalf("failed to list pods: %v", err)
	}

	if len(pods) != 1 {
		t.Errorf("expected 1 pod, got %d", len(pods))
	}

	pod.Status = types.PodStatusRunning
	err = store.UpdatePod(pod)
	if err != nil {
		t.Fatalf("failed to update pod: %v", err)
	}

	retrieved, _ = store.GetPod("default", "test-pod")
	if retrieved.Status != types.PodStatusRunning {
		t.Errorf("expected status 'Running', got '%s'", retrieved.Status)
	}

	err = store.DeletePod("default", "test-pod")
	if err != nil {
		t.Errorf("failed to delete pod: %v", err)
	}

	pods, _ = store.ListPods()
	if len(pods) != 0 {
		t.Errorf("expected 0 pods after deletion, got %d", len(pods))
	}

	t.Log("all store tests passed!")
}
