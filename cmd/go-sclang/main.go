package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	models "github.com/WnP/go-sclang/pkg"
)

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "host", "localhost", "Server host")
	flag.StringVar(&host, "h", "localhost", "-host shorthand ")

	flag.IntVar(&port, "port", 5533, "Server port")
	flag.IntVar(&port, "p", 5533, "-port shorthand ")

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage of %s:
Wrap sclang CLI in order to provide an HTTP API, expose:

POST application/json {"Code": "'sclang code'.postln, "Stdout": true}

if Stdout is true then sclang return value is returned.

Flags:

`,
			os.Args[0],
		)

		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {

	input := make(chan models.Query)
	output := make(chan string)
	requireOutput := make(chan interface{})

	cmd := exec.Command("sclang")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Cannot pipe stdin: %v\n", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Cannot pipe stdout: %v\n", err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Cannot pipe stderr: %v\n", err.Error())
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("Cannot start sclang: %v\n", err.Error())
	}

	closeIn := make(chan interface{}, 1)
	closeHandler(closeIn)
	go handleStdin(stdin, closeIn, input, requireOutput)
	go handleStdout(stdout, requireOutput, output)
	go handleStderr(stderr)

	go startHTTPServer(cmd, input, closeIn, output)

	if err := cmd.Wait(); err != nil {
		log.Fatalf("sclang fail: %v\n", err.Error())
	}
}

// Shutdown gracefully
func closeHandler(closeIn chan<- interface{}) {
	var i interface{}
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		closeIn <- i
	}()

}

// Handle HTTP POST requests
func getHTTPHandler(
	cmd *exec.Cmd,
	input chan<- models.Query,
	closeIn chan<- interface{},
	output <-chan string,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request\n")
			return
		}

		var query models.Query
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&query); err != nil {
			log.Printf("Bad payload: %v", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request\n")
			return
		}

		input <- query

		if query.Stdout {
			out := <-output
			fmt.Fprintf(w, out)
		}

		if query.Kill || query.Reload {
			input <- models.Query{
				Code:   "Server.quitAll;",
				Stdout: true, // in order to wait for kill to be done
			}
			<-output
		}
		if query.Kill {
			var i interface{}
			log.Println("Killing sclang...")
			closeIn <- i
		} else if query.Reload {
			log.Println("Reloading sclang...")
			if err := cmd.Process.Signal(syscall.SIGUSR1); err != nil {
				log.Fatalf("Cannot reload sclang: %v", err.Error())
			}
			// Send an empty string in order to trigger sclang handleSigUsr1()
			input <- models.Query{Code: ""}
		}
	}
}

func startHTTPServer(
	cmd *exec.Cmd,
	input chan<- models.Query,
	closeIn chan<- interface{},
	output <-chan string,
) {
	http.HandleFunc("/", getHTTPHandler(cmd, input, closeIn, output))
	host := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Serving go-sclang server on %s\n", host)
	http.ListenAndServe(host, nil)
}

func handleStdin(stdin io.WriteCloser, sigTerm <-chan interface{}, input <-chan models.Query, requireOutput chan<- interface{}) {
	defer stdin.Close()
	var i interface{}
	for {
		select {
		case <-sigTerm:
			return
		case in := <-input:
			/*
				0x0c :executes the current command line as SC code and prints the result to stdout
				0x1b :executes the currently accumulated command line as SC code
			*/
			var e rune
			if in.Stdout {
				e = '\x0c'
				go func() { requireOutput <- i }()
			} else {
				e = '\x1b'
			}
			stdin.Write([]byte(in.Code + string(e)))
		}
	}
}

func handleStdout(stdout io.ReadCloser, requireOutput <-chan interface{}, output chan<- string) {
	var send bool = false // This is the default value but... explicit is better than implicit
	for {
		val := make([]byte, 1024)
		if _, err := stdout.Read(val); err != nil {
			log.Fatalf("Cannot read stdout: %v\n", err.Error())
		} else {
			select {
			case <-requireOutput:
				send = true
			default:
			}
			v := string(val)
			fmt.Fprint(os.Stdout, v)
			if send {
				for _, line := range strings.Split(v, "\n") {
					if strings.HasPrefix(line, "-> ") {
						output <- line[3:]
						send = false
					}
				}
			}

		}
	}
}

func handleStderr(stderr io.ReadCloser) {
	for {
		val := make([]byte, 1024)
		if _, err := stderr.Read(val); err != nil {
			log.Fatalf("Cannot read stderr: %v\n", err.Error())
		} else {
			fmt.Fprint(os.Stderr, string(val))
		}
	}
}
