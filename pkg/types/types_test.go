package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPodJSON(t *testing.T) {
	pod := Pod{
		Name:      "nginx-pod",
		Namespace: "default",
		Image:     "nginx:latest",
		Status:    PodStatusPending,
		CreatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(pod)
	if err != nil {
		t.Fatalf("Failed to marshal pod: %v", err)
	}

	t.Logf("Pod JSON: %s", string(jsonData))

	var pod2 Pod
	err = json.Unmarshal(jsonData, &pod2)
	if err != nil {
		t.Fatalf("Failed to unmarshal pod: %v", err)
	}

	if pod2.Name != "nginx-pod" {
		t.Errorf("Expected name 'nginx-pod', got '%s'", pod2.Name)
	}
}
