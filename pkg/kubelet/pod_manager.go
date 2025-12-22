package kubelet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vladrosant/m8s/pkg/types"
)

type PodManager struct {
	nodeName string
	apiURL   string
	runtime  *ContainerRuntime
}

func NewPodManager(nodeName, apiURL string) *PodManager {
	return &PodManager{
		nodeName: nodeName,
		apiURL:   apiURL,
		runtime:  NewContainerRuntime(),
	}
}

func (pm *PodManager) Sync() error {
	desiredPods, err := pm.getPodsFromAPI()
	if err != nil {
		return fmt.Errorf("failed to get pods from API: %w", err)
	}

	runningContainers, err := pm.runtime.ListRunningContainers()
	if err != nil {
		return fmt.Errorf("failed to list running containers: %w", err)
	}

	for _, pod := range desiredPods {
		if pod.NodeName != pm.nodeName {
			continue
		}

		containerName := fmt.Sprintf("m8s-%s-%s", pod.Namespace, pod.Name)

		isRunning := false
		for _, name := range runningContainers {
			if name == containerName {
				isRunning = true
				break
			}
		}

		if !isRunning {
			fmt.Printf("Starting pod %s/%s...\n", pod.Namespace, pod.Name)
			if err := pm.runtime.StartContainer(pod); err != nil {
				fmt.Printf("Failed to start pod %s/%s: %v\n", pod.Namespace, pod.Name, err)
				pm.updatePodStatus(pod, types.PodStatusFailed)
				continue
			}
		}

		status, err := pm.runtime.GetContainerStatus(pod)
		if err != nil {
			fmt.Printf("Failed to get status for pod %s/%s: %v\n", pod.Namespace, pod.Name, err)
			continue
		}

		if pod.Status != status {
			pm.updatePodStatus(pod, status)
		}
	}

	desiredContainers := make(map[string]bool)
	for _, pod := range desiredPods {
		if pod.NodeName == pm.nodeName {
			containerName := fmt.Sprintf("m8s-%s-%s", pod.Namespace, pod.Name)
			desiredContainers[containerName] = true
		}
	}

	for _, containerName := range runningContainers {
		if !desiredContainers[containerName] {
			fmt.Printf("Stopping unexpected container %s...\n", containerName)

			parts := strings.Split(containerName, "-")
			if len(parts) >= 3 {
				pod := types.Pod{
					Namespace: parts[1],
					Name:      strings.Join(parts[2:], "-"),
				}
				pm.runtime.StopContainer(pod)
			}
		}
	}

	return nil
}

func (pm *PodManager) getPodsFromAPI() ([]types.Pod, error) {
	url := fmt.Sprintf("%s/api/v1/pods", pm.apiURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var podList types.PodList
	if err := json.NewDecoder(resp.Body).Decode(&podList); err != nil {
		return nil, err
	}

	return podList.Items, nil
}

func (pm *PodManager) updatePodStatus(pod types.Pod, status string) error {
	fmt.Printf("Pod %s/%s status: %s -> %s\n", pod.Namespace, pod.Name, pod.Status, status)
	return nil
}
