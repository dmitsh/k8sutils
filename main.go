package main

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/dmitsh/k8sutils/pkg/kubeclient"
)

var client *kubeclient.Client

func getRoot(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	//fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "This is my website!\n")
}

func doPodCopy(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//w.Header().Set("Content-Type", "application/json")

	//fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello, HTTP!\n")
	inout, outout, errout, err := client.PodCopyFile("default", "/tmp/my.txt", "/tmp/copy-my.txt", "mypod")

	log.Infof("INOUT: %s", inout.String())
	log.Infof("OUTOUT: %s", outout.String())
	log.Infof("ERROUT: %s", errout.String())

	// Write response body
	//w.Write([]byte("Hello, World!"))

	// Write response status code
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	var err error
	client, err = kubeclient.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/cp", doPodCopy)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("DONE")
	}
}
