package kubelet

import (
	"fmt"
	"os/exec"
	"strings"

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

func (cr *ContainerRuntime) StopContainer(pod types.Pod) error {
	containerName := fmt.Sprintf("m8s-%s-%s", pod.Namespace, pod.Name)

	exists, err := cr.containerExists(containerName)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	cmd := exec.Command("docker", "stop", containerName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove container: %w, output: %s", err, string(output))
	}

	cmd = exec.Command("docker", "rm", containerName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove container: %w, output: %s", err, string(output))
	}

	return nil
}

func (cr *ContainerRuntime) GetContainerStatus(pod types.Pod) (string, error) {
	containerName := fmt.Sprintf("m8s-%s-%s", pod.Namespace, pod.Name)

	exists, err := cr.containerExists(containerName)
	if err != nil {
		return types.PodStatusFailed, err
	}

	if !exists {
		return types.PodStatusPending, nil
	}

	running, err := cr.isContainerRunning(containerName)
	if err != nil {
		return types.PodStatusFailed, err
	}

	if running {
		return types.PodStatusRunning, nil
	}

	exited, err := cr.hasContainerExited(containerName)
	if err != nil {
		return types.PodStatusFailed, err
	}

	if exited {
		return types.PodStatusSucceeded, nil
	}

	return types.PodStatusFailed, nil
}

func (cr *ContainerRuntime) containerExists(name string) (bool, error) {
	cmd := exec.Command("docker", "ps", "a","--filter", fmt.Sprintf("name=^%s$", name), "format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == name, nil
}

func (cr *ContainerRuntime) isContainerRunning(name string) (bool, error) {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=^%s$", name), "format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == name, nil
}

func (cr *ContainerRuntime) hasContainerExited(name string) (bool, error) {
	cmd := exec.Command("docker", "inspect", "--format", "{{.State.ExitCode}}", name")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
}
