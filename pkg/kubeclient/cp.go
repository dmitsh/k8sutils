package kubeclient

import (
	"bytes"
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/cp"
	"k8s.io/kubectl/pkg/cmd/util"
)

func (c *Client) PodCopyFile(namespace, src, dst, container string) (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
	ioStreams, in, out, errOut := genericclioptions.NewTestIOStreams()
	opt := cp.NewCopyOptions(ioStreams)
	opt.Clientset = c.client
	opt.ClientConfig = c.config
	opt.Container = container
	opt.Namespace = namespace

	if err := opt.Validate(); err != nil {
		return nil, nil, nil, fmt.Errorf("could not validate copy options: %v", err)
	}

	nf := util.NewFactory(&genericclioptions.ConfigFlags{})
	cmd := cp.NewCmdCp(nf, ioStreams)

	if err := opt.Complete(nf, cmd, []string{src, dst}); err != nil {
		return nil, nil, nil, fmt.Errorf("could not complete copy options: %v", err)
	}

	if err := opt.Run(); err != nil {
		return nil, nil, nil, fmt.Errorf("could not copy: %v", err)
	}

	return in, out, errOut, nil
}
