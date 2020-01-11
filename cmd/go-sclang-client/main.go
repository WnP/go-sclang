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

	models "github.com/WnP/go-sclang/pkg"
)

var (
	host   string
	port   int
	stdout bool
	kill   bool
	reload bool
)

func init() {
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

	flag.Parse()
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

	send(query)
}

func send(query models.Query) {

	b, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Cannot marshal json: %s\n", err.Error())
	}

	resp, err := http.Post(
		fmt.Sprintf("%s:%d/", host, port),
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		log.Fatalf("HTTP query failed: %s\n", err.Error())
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
