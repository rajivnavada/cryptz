package main

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"flag"
	"fmt"
	pb "github.com/rajivnavada/cryptz_pb"
	"io"
	"os"
	"strings"
)

var (
	host  = flag.String("host", "127.0.0.1", "HTTP service host")
	port  = flag.String("port", "8000", "HTTP port at which the service will run")
	debug = flag.Bool("debug", false, "Turn on debug mode")
	fpr   = flag.String("fpr", "", "Fingerprint of key to use")
)

func repl(cli Client) {
	bio := bufio.NewReader(os.Stdin)
	for {
		// Scan the line from STDIN
		line, err := bio.ReadString('\n')
		if err != nil {
			logError(err, "Error reading line")
			continue
		}

		tokens := strings.SplitN(strings.TrimSpace(line), " ", 5)
		// Normalize to at least 5 tokens
		for i := len(tokens); i < 5; i++ {
			tokens = append(tokens, "")
		}

		// Interpret the printed line
		op := &pb.Operation{}

		// Prepare operation
		switch tokens[0] {
		case "project":
			o := &pb.ProjectOperation{}
			switch tokens[1] {
			case "list":
				o.Command = pb.ProjectOperation_LIST
				break

			case "create":
				o.Command = pb.ProjectOperation_CREATE
				o.Name = tokens[2]
				o.Environment = tokens[3]
				break

			case "delete":
				//o.ProjectId = a3
				break
			}
			op.ProjectOrCredentialOp = &pb.Operation_ProjectOp{ProjectOp: o}
			break

		case "credential":
			//op = &pb.Operation{ProjectOrCredentialOp: &Operation_CredentialOp{CredentialOp: o}}
			break

		case "quit":
			cli.Close()
			return
		}

		// Write to client
		cli.Send(op)
	}
}

func main() {
	flag.Parse()

	// TODO validate args

	// Read cert.pem and key.pem into a buffer
	buf := &bytes.Buffer{}
	for _, fname := range []string{"cert.pem", "key.pem"} {
		if f, err := os.Open(fname); err != nil {
			panic(err)
		} else {
			io.Copy(buf, f)
		}
	}

	// Create a cert pool
	certs := x509.NewCertPool()
	if !certs.AppendCertsFromPEM(buf.Bytes()) {
		println("Could not parse cert from PEM")
		return
	}

	wssurl := fmt.Sprintf("wss://%s:%s/ws/%s", *host, *port, *fpr)
	origin := fmt.Sprintf("https://%s:%s", *host, *port)

	// Start a websocket client to receive/send messages
	cli := NewClient(wssurl, origin, certs)

	// Run the client
	go cli.Run()

	// Start REPL with client
	repl(cli)
}
