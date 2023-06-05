package kubeclient

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

func (client *Client) ExecCmd(ctx context.Context, namespace, pod, container string, cmd []string, stdout io.Writer, stderr io.Writer) error {
	req := client.client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   cmd,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(client.config, "POST", req.URL())
	if err != nil {
		log.Error("failed on NewSPDYExecutor")
		return err
	}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: stdout,
		Stderr: stderr,
	})

	if err != nil {
		log.Error("failed on StreamWithContext")
		return err
	}

	return nil
}
