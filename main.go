package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/dmitsh/k8sutils/pkg/kubeclient"
)

var (
	client                    *kubeclient.Client
	namespace, pod, container string
	port                      int
)

func doUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	io.WriteString(w, fmt.Sprintf("Usage: curl %v/exec\n", ctx.Value("serverAddr")))
}

func doExec(w http.ResponseWriter, r *http.Request) {
	var stdout, stderr bytes.Buffer
	err := client.ExecCmd(context.Background(), namespace, pod, container,
		[]string{"cat", "/etc/test/result.json"}, &stdout, &stderr)

	log.Infof("STDOUT: %s", stdout.String())
	log.Infof("STDERR: %s", stderr.String())

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(stdout.Bytes())
		w.WriteHeader(http.StatusOK)
	} else {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	flag.StringVar(&namespace, "namespace", "default", "namespace")
	flag.StringVar(&pod, "pod", "test", "pod name")
	flag.StringVar(&container, "container", "test", "container name")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	var err error
	client, err = kubeclient.New()
	if err != nil {
		log.Fatalf("failed to create Kubernetes client: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", doUsage)
	mux.HandleFunc("/exec", doExec)

	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal(err.Error())
	}
}
