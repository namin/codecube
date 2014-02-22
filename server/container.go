package main

import (
	"log"
	"time"
	"io"
	"bufio"
	dcli "github.com/fsouza/go-dockerclient"
)

func tailOutput(name string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log.Printf("[%s] %s\n", name, scanner.Text())
	}
}

func notmain() {
	client, err := dcli.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		panic(err)
	}

	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()
	runner := NewRunner(client, "ruby", "puts \"yo i'm rubby #{7*7}\"")
	runner.OutStream = outWriter
	runner.ErrStream = errWriter

	go tailOutput("stdout", outReader)
	go tailOutput("stderr", errReader)

	log.Println("Running code...")
	if _, err := runner.Run(10000); err != nil {
	    panic(err)
	}

	outReader.Close()
	errReader.Close()

	time.Sleep(1e9)
}

