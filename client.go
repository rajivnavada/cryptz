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
	"time"
)

const (
	writeWait = 10 * time.Second
)

func logError(err error, info string) {
	println(info)
	if err != nil {
		println(err)
	}
	println("")
}

type Client interface {
	Send(op *pb.Operation) error
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
			if response.Status == pb.Response_ERROR {
				logError(fmt.Errorf(response.Error), "Error when creating project")
				continue
			}
			if response.Status == pb.Response_SUCCESS {

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
					}

				} else if credResponse != nil {

				}

			}
			// Display results

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

func (c *client) Send(op *pb.Operation) error {
	message, err := proto.Marshal(op)
	if err != nil {
		return err
	}
	c.WriteChan <- message
	return nil
}

func (c *client) Run() {
	dialer := &websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            c.CertPool,
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
