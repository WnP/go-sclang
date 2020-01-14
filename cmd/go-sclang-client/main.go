package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	models "github.com/WnP/go-sclang/pkg"
)

var (
	host         string
	port         int
	stdout       bool
	kill         bool
	reload       bool
	retry        bool
	timeout      time.Duration
	retryTimeout time.Duration
	client       http.Client
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage of %s:
Reads code from stdin and send it to go-sclang server, eg:

	echo '"Hello from %s".postln' | %s

Flags:

`,
			os.Args[0],
			os.Args[0],
			os.Args[0],
		)

		flag.PrintDefaults()
	}

	flag.StringVar(&host, "host", "http://localhost", "Server host")
	flag.StringVar(&host, "h", "http://localhost", "-host shorthand ")

	flag.IntVar(&port, "port", 5533, "Server port")
	flag.IntVar(&port, "p", 5533, "-port shorthand ")

	flag.BoolVar(&stdout, "stdout", false, "If true return sclang value")
	flag.BoolVar(&stdout, "o", false, "-stdout shorthand ")

	flag.BoolVar(&kill, "kill", false, "If true kill sclang server")
	flag.BoolVar(&kill, "k", false, "-kill shorthand ")

	flag.BoolVar(&reload, "reload", false, "If true reload sclang server by sending SIGUSR1")
	flag.BoolVar(&reload, "r", false, "-kill shorthand ")

	flag.BoolVar(&retry, "retry", false, "If true will retry if query fail until timeout")

	var t string
	var tr string
	// See https://golang.org/pkg/time/#example_ParseDuration for availabe units
	flag.StringVar(&t, "timeout", "10s", "sclang server timeout")
	flag.StringVar(&t, "t", "10s", "-timeout shorthand")

	flag.StringVar(&tr, "retry-timeout", "5s", "Retry timeout")

	flag.Parse()
	var err error
	if timeout, err = time.ParseDuration(t); err != nil {
		log.Fatalf("Invalid timeout: %v\n", err.Error())
	}
	if retryTimeout, err = time.ParseDuration(tr); err != nil {
		log.Fatalf("Invalid retry timeout: %v\n", err.Error())
	}

	client = http.Client{
		Timeout: timeout,
	}
}

func main() {
	var code []byte
	var err error

	if !kill && !reload {
		if code, err = ioutil.ReadAll(os.Stdin); err != nil {
			log.Fatalf("Cannot read from stdin: %s\n", err.Error())
		}
	}

	query := models.Query{
		Code:   string(code),
		Stdout: stdout,
		Kill:   kill,
		Reload: reload,
	}

	deadline := time.Now().Add(retryTimeout)
	send(query, deadline)
}

func send(query models.Query, deadline time.Time) {

	b, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Cannot marshal json: %s\n", err.Error())
	}

	resp, err := client.Post(
		fmt.Sprintf("%s:%d/", host, port),
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		if retry && time.Now().Before(deadline) {
			time.Sleep(10 * time.Millisecond)
			send(query, deadline)
		} else {
			log.Fatalf("HTTP query failed: %s\n", err.Error())
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf(
			"Wrong returned status: (%d) %s\n",
			resp.StatusCode, resp.Status,
		)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Cannot read response body: %s\n", err.Error())
		}
		log.Fatalf("%s\n", body)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Cannot read response body: %s\n", err.Error())
	}
	fmt.Printf("%s", body)
}
