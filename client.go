package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	pb "github.com/rajivnavada/cryptz_pb"
	"github.com/rajivnavada/gpgme"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	writeWait = 10 * time.Second
)

func logError(err error, info string) {
	println(info)
	if err != nil {
		println(err.Error())
	}
	println("")
}

type Client interface {
	Send(op *pb.Operation) (<-chan bool, error)
	Run()
	Close()
}

type client struct {
	URL       string
	Origin    string
	CertPool  *x509.CertPool
	WriteChan chan []byte
	quitChan  chan bool
}

var (
	requestIdLock       = &sync.Mutex{}
	requestId     int32 = 0
	openRequests        = make(map[int32]chan<- bool)
)

func openRequest() (int32, chan bool) {
	requestIdLock.Lock()
	defer requestIdLock.Unlock()

	requestId++
	ch := make(chan bool)
	openRequests[requestId] = ch
	return requestId, ch
}

func finalizeRequest(reqId int32) {
	requestIdLock.Lock()
	defer requestIdLock.Unlock()

	if ch, ok := openRequests[requestId]; !ok {
		return
	} else {
		ch <- true
		close(ch)
		delete(openRequests, requestId)
	}
}

func (c *client) readPump(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			logError(nil, "Connection closed. Bye!")
			return
		}
		// Handle text messages
		switch messageType {
		case websocket.TextMessage:
			// decrypt the message before displaying
			if result, err := gpgme.DecryptMessage(string(p)); err != nil {
				println("------------------------------------------------------------")
				println("An error occured when trying to decrypt message")
				println("Ignoring this message")
				println("------------------------------------------------------------")
			} else {
				println("------------------------------------------------------------")
				println(result)
				println("------------------------------------------------------------")
			}
			break

		case websocket.BinaryMessage:
			response := &pb.Response{}
			// Unmarshal the message using protocol buffers
			err := proto.Unmarshal(p, response)
			if err != nil {
				logError(err, "Could not unmarshal response from server")
				continue
			}

			// Decrypt relevant parts
			switch response.Status {
			case pb.Response_ERROR:
				logError(fmt.Errorf(response.Error), "Server responded with an error")

			case pb.Response_SUCCESS:
				projResponse := response.GetProjectOpResponse()
				credResponse := response.GetCredentialOpResponse()

				if projResponse != nil {

					switch projResponse.Command {
					case pb.ProjectOperation_CREATE:
						fmt.Printf("%s\n", response.Info)

					case pb.ProjectOperation_UPDATE:
						fmt.Printf("%s\n", response.Info)

					case pb.ProjectOperation_DELETE:
						fmt.Printf("%s\n", response.Info)

					case pb.ProjectOperation_LIST:

					case pb.ProjectOperation_LIST_CREDENTIALS:
						fmt.Printf("%s\n", response.Info)
						for _, c := range projResponse.Credentials {
							fmt.Printf("%s (Id = %d)\n", c.Key, c.Id)
						}

					case pb.ProjectOperation_ADD_MEMBER:
						fmt.Printf("%s\n", response.Info)

					case pb.ProjectOperation_DELETE_MEMBER:
						fmt.Printf("%s\n", response.Info)
					}

				} else if credResponse != nil {

					switch credResponse.Command {
					case pb.CredentialOperation_GET:
						cred := credResponse.GetCredential()
						if cred == nil || cred.Cipher == "" {
							logError(nil, "Server returned an empty response for credential request.")
						} else {
							if value, err := gpgme.DecryptMessage(credResponse.Credential.Cipher); err != nil {
								logError(err, "Could not decrypt credential cipher")
							} else {
								fmt.Printf("%s\n", value)
							}
						}

					case pb.CredentialOperation_SET:
						fmt.Printf("%s\n", response.Info)

					case pb.CredentialOperation_DELETE:
						fmt.Printf("%s\n", response.Info)
					}

				}
			}
			reqId := response.OpId
			finalizeRequest(reqId)

		case websocket.CloseMessage:
			return
		}
	}
}

func (c *client) writePump(conn *websocket.Conn) {
	for {
		select {
		case message := <-c.WriteChan:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil && websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logError(nil, "Connection closed. Bye!")
				return
			}
			if err != nil {
				logError(err, "An error occured when trying to perform operation")
			}

		case <-c.quitChan:
			logError(nil, "Connection closed. Bye!")
			return
		}
	}
}

func (c *client) Send(op *pb.Operation) (<-chan bool, error) {
	// Add the ID to the message
	reqId, out := openRequest()
	op.OpId = reqId

	message, err := proto.Marshal(op)
	if err != nil {
		finalizeRequest(reqId)
		return nil, err
	}

	c.WriteChan <- message
	return out, nil
}

func (c *client) Run() {
	dialer := &websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	header := http.Header{
		"Origin": {c.Origin},
	}

	conn, _, err := dialer.Dial(c.URL, header)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go c.writePump(conn)
	c.readPump(conn)

	c.Close()
}

func (c *client) Close() {
	c.quitChan <- true
	close(c.WriteChan)
	close(c.quitChan)
	os.Exit(0)
}

func NewClient(url, origin string, pool *x509.CertPool) Client {
	return &client{
		URL:       url,
		Origin:    origin,
		CertPool:  pool,
		WriteChan: make(chan []byte),
		quitChan:  make(chan bool),
	}
}
