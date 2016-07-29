package main

import (
	"bufio"
	"flag"
	"fmt"
	pb "github.com/rajivnavada/cryptz_pb"
	"io"
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
		if err != nil && err != io.EOF {
			logError(err, "Error reading line")
			continue
		}
		if err == io.EOF {
			println("Closing connection. Bye!")
			cli.Close()
			return
		}

		tokens := strings.SplitN(strings.TrimSpace(line), " ", 6)
		// Normalize to at least 6 tokens
		for i := len(tokens); i < 6; i++ {
			tokens = append(tokens, "")
		}

		// Operations:
		// list projects
		// pin project ID
		// project create NAME ENVIRONMENT
		// project ID list credentials
		// project ID list members
		// project ID add credential CREDKEY
		// project ID remove credential CREDKEY
		// project ID add member MEMBEREMAIL
		// project ID remove member MEMBEREMAIL
		// project ID remove
		// unpin
		// list messages
		// message ID show

		// Interpret the printed line
		op := &pb.Operation{}

		switch tokens[0] {
		case "":
			continue

		case "list":
			switch tokens[1] {
			case "projects":
				o := &pb.ProjectOperation{}
				o.Command = pb.ProjectOperation_LIST

			case "messages":
				continue

			default:
				continue
			}

		case "project":
			o := &pb.ProjectOperation{}
			op.ProjectOrCredentialOp = &pb.Operation_ProjectOp{ProjectOp: o}

			switch tokens[1] {
			case "create":
				o.Command = pb.ProjectOperation_CREATE
				o.Name = tokens[2]
				o.Environment = tokens[3]

			default:
				pid, err := strconv.Atoi(tokens[1])
				if err != nil {
					logError(err, "Could not convert project id to int")
					continue
				}
				o.ProjectId = int32(pid)

				switch tokens[2] {
				case "list":
					switch tokens[3] {
					case "credentials":
						o.Command = pb.ProjectOperation_LIST_CREDENTIALS

					case "members":
						continue
					}

				case "add":
					switch tokens[3] {
					case "member":
						o.Command = pb.ProjectOperation_ADD_MEMBER
						o.MemberEmail = tokens[4]

					case "credential":
						o.Command = pb.CredentialOperation_SET
						o.Key = tokens[4]
						o.Value = tokens[5]
					}

				case "remove":
					switch tokens[3] {
					case "member":
						o.Command = pb.ProjectOperation_DELETE_MEMBER
						mid, err := strconv.Atoi(tokens[4])
						if err != nil {
							logError(err, "Could not convert member id to int")
							continue
						}
						o.MemberId = int32(mid)

					case "credential":
						o.Command = pb.CredentialOperation_DELETE
						o.Key = tokens[4]

					default:
						o.Command = pb.ProjectOperation_DELETE
					}

				case "get":
					switch tokens[3] {
					case "credential":
						o.Command = pb.CredentialOperation_GET
						o.Key = tokens[4]

					default:
						continue
					}
				}
			}

		case "message":
			continue
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
