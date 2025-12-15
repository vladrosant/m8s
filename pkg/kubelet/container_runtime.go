package kubelet

import (
	"fmt"
	"os/exec"

	"github.com/vladrosant/m8s/pkg/types"
)

type ContainerRuntime struct{}

func NewContainerRuntime() *ContainerRuntime {
	return &ContainerRuntime{}
}

func (cr *ContainerRuntime) StartContainer(pod types.Pod) error {
	containerName := fmt.Sprintf("m8s-%s-%s", pod.Namespace, pod.Name)

	exists, err := cr.containerExists(containerName)
	if err != nil {
		return fmt.Errorf("failed to check if container exists: %w", err)
	}

	if exists {
		running, err := cr.isContainerRunning(containerName)
		if err != nil {
			return fmt.Errorf("failed to check container status: %w", err)
		}

		if running {
			return nil
		}

		return cr.startExistingContainer(containerName)

	}

	cmd := exec.Command("docker", "run", "-d", "--name", containerName, pod.Image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %w, output: %s", err, string(output))
	}

	return nil
}
