package main

import (
	"bufio"
	"flag"
	"fmt"
	pb "github.com/rajivnavada/cryptz_pb"
	"os"
	"strconv"
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
		fmt.Print("> ")
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

			case "list-credentials":
				o.Command = pb.ProjectOperation_LIST_CREDENTIALS
				pid, err := strconv.Atoi(tokens[2])
				if err != nil {
					logError(err, "Could not convert project id to int")
					continue
				}
				o.ProjectId = int32(pid)

			case "create":
				o.Command = pb.ProjectOperation_CREATE
				o.Name = tokens[2]
				o.Environment = tokens[3]

			case "delete":
				o.Command = pb.ProjectOperation_DELETE
				pid, err := strconv.Atoi(tokens[2])
				if err != nil {
					logError(err, "Could not convert project id to int")
					continue
				}
				o.ProjectId = int32(pid)
			}
			op.ProjectOrCredentialOp = &pb.Operation_ProjectOp{ProjectOp: o}

		case "member":
			o := &pb.ProjectOperation{}
			switch tokens[1] {
			case "add":
				o.Command = pb.ProjectOperation_ADD_MEMBER
				pid, err := strconv.Atoi(tokens[2])
				if err != nil {
					logError(err, "Could not convert project id to int")
					continue
				}
				o.ProjectId = int32(pid)
				o.MemberEmail = tokens[3]

			case "delete":
				o.Command = pb.ProjectOperation_DELETE_MEMBER
				mid, err := strconv.Atoi(tokens[2])
				if err != nil {
					logError(err, "Could not convert member id to int")
					continue
				}
				o.MemberId = int32(mid)
			}
			op.ProjectOrCredentialOp = &pb.Operation_ProjectOp{ProjectOp: o}

		case "credential":
			o := &pb.CredentialOperation{}
			switch tokens[1] {
			case "set":
				o.Command = pb.CredentialOperation_SET
				pid, err := strconv.Atoi(tokens[2])
				if err != nil {
					logError(err, "Could not convert project id to integer")
					continue
				}
				o.Project = int32(pid)
				o.Key = tokens[3]
				o.Value = tokens[4]

			case "get":
				o.Command = pb.CredentialOperation_GET
				pid, err := strconv.Atoi(tokens[2])
				if err != nil {
					logError(err, "Could not convert project id to integer")
					continue
				}
				o.Project = int32(pid)
				o.Key = tokens[3]

			case "delete":
				o.Command = pb.CredentialOperation_DELETE
				pid, err := strconv.Atoi(tokens[2])
				if err != nil {
					logError(err, "Could not convert project id to integer")
					continue
				}
				o.Project = int32(pid)
				o.Key = tokens[3]
			}
			op.ProjectOrCredentialOp = &pb.Operation_CredentialOp{CredentialOp: o}

		case "quit":
			cli.Close()
			return
		}

		// Write to client
		if waitCh, err := cli.Send(op); err != nil {
			logError(err, "Error when trying to send operation")
		} else {
			<-waitCh
		}
	}
}

func main() {
	flag.Parse()

	// TODO validate args

	fingerprint := strings.Replace(*fpr, " ", "", -1)
	if fingerprint == "" {
		println("A key fingerprint is required to run the client")
		return
	}

	wssurl := fmt.Sprintf("wss://%s:%s/ws/%s", *host, *port, fingerprint)
	origin := fmt.Sprintf("https://%s:%s", *host, *port)

	// Start a websocket client to receive/send messages
	cli := NewClient(wssurl, origin, nil)

	// Run the client
	go cli.Run()

	// Start REPL with client
	repl(cli)
}
