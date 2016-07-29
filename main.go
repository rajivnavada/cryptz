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

func handleError(err error) {
	if err != nil {
		logError(err, "Could not create command")
	}
}

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
		var op *pb.Operation

		switch tokens[0] {
		case "":
			continue

		case "list":
			switch tokens[1] {
			case "projects":
				op = NewProjectListOperation()

			case "messages":
				continue

			default:
				continue
			}

		case "project":
			switch tokens[1] {
			case "create":
				op, err = NewProjectCreateOperation(tokens[3], tokens[4])
				handleError(err)

			default:
				pid, err := strconv.Atoi(tokens[1])
				if err != nil {
					handleError(err)
					continue
				}
				p := Project(pid)

				switch tokens[2] {
				case "list":
					switch tokens[3] {
					case "credentials":
						op, err = p.NewListCredentialsOperation()
						handleError(err)

					case "members":
						continue
					}

				case "add":
					switch tokens[3] {
					case "member":
						op, err = p.NewAddMemberOperation(tokens[4])
						handleError(err)

					case "credential":
						op, err = p.NewAddCredentialOperation(tokens[4], tokens[5])
						handleError(err)
					}

				case "remove":
					switch tokens[3] {
					case "":
						op, err = p.NewDeleteOperation()
						handleError(err)

					case "member":
						mid, err := strconv.Atoi(tokens[4])
						if err != nil {
							handleError(err)
							continue
						}
						op, err = p.NewDeleteMemberOperation(mid)
						handleError(err)

					case "credential":
						op, err = p.NewDeleteCredentialOperation(tokens[4])
						handleError(err)

					default:
						continue
					}

				case "get":
					switch tokens[3] {
					case "credential":
						op, err = p.NewGetCredentialOperation(tokens[4])
						handleError(err)

					default:
						continue
					}
				}
			}

		case "message":
			continue

		default:
			continue
		}

		if err != nil || op == nil {
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
