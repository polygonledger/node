package main

//basic server based protocol
//server receives tx messages
//adds tx messages to a pool
//every 10 tx messages a new block gets created

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"encoding/gob"
	"net/http"

	block "github.com/polygon/block"

	protocol "github.com/polygon/net"
)

//var latestBlock block.Block
var blockheight int = 0
var tx_pool []block.Tx

//TODO
//handleNewTx => put in txpool
//createNewBlock
//newWallet

/*
Outgoing connections
 A `net.Conn` satisfies the io.Reader and `io.Writer` interfaces
*/

// connects to a TCP Address and returns a connection with a timeout and wrapped into a buffered ReadWriter
func Open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	// Note that the local port is chosen on the fly. If the local port
	// must be a specific one, use DialTCP() instead.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

/*
Incoming connections
receive the name of a command terminated by `\n`, followed by data.
create an `Endpoint` object with the following properties:
* It allows to register one or more handler functions, where each can handle a
  particular command.
* It dispatches incoming commands to the associated handler based on the commands
  name.
*/

// HandleFunc is a function that handles an incoming command
// It receives the open connection wrapped in a `ReadWriter` interface
type HandleFunc func(*bufio.ReadWriter)

// Endpoint provides an endpoint to other processess
// that they can send data to.
type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc

	// Maps are not threadsafe, so we need a mutex to control access
	m sync.RWMutex
}

// NewEndpoint creates a new endpoint. To keep things simple,
// the endpoint listens on a fixed port number
func NewEndpoint() *Endpoint {
	// Create a new Endpoint with an empty list of handler funcs
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

// AddHandleFunc adds a new function for handling incoming data
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

// Listen starts listening on the endpoint port on all interfaces

func (e *Endpoint) Listen() error {
	var err error
	e.listener, err = net.Listen("tcp", protocol.Port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", protocol.Port)
	}
	log.Println("Listen on", e.listener.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go e.handleMessages(conn)
	}
}

// handleMessages reads the connection up to the first newline
// Based on this string, it calls the appropriate HandleFunc
func (e *Endpoint) handleMessages(conn net.Conn) {
	// Wrap the connection into a buffered reader for easier reading.
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// Read from the connection until EOF. Expect a command name as the
	// next input. Call the handler that is registered for this command.
	for {
		log.Print("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}
		// Trim the request string - ReadString does not strip any newlines.
		cmd = strings.Trim(cmd, "\n ")
		log.Println(cmd)

		// Fetch the appropriate handler function from the 'handler' map and call it.
		e.m.RLock()
		handleCommand, ok := e.handler[cmd]
		e.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handleCommand(rw)
	}
}

// handleGob handles "GOB" request. decodes the received GOB data into a struct
func handleGob(rw *bufio.ReadWriter) {
	log.Print("Receive GOB data:")
	var data block.Tx
	// Create a decoder that decodes directly into a struct variable.
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	log.Printf("data: \n%#v\n", data.Nonce)

	timestamp := time.Now().Unix()
	b := []byte(string(timestamp))
	log.Println(b)
	hash := sha256.Sum256(b)
	//log.Print(fmt.Sprintf("%x", )[:45])
	//hash := sha256.Sum256(xdata)

	data.Id = hash

	tx_pool = append(tx_pool, data)
	//blockheight += 1

	//log.Printf("tx_pool_size: \n%d", tx_pool_size)

	log.Printf("tx_pool size: \n%d", len(tx_pool))
}

// server listens for incoming requests and dispatches them to
// registered handler functions
func server() error {
	endpoint := NewEndpoint()

	// Add the handle funcs.
	endpoint.AddHandleFunc("GOB", handleGob)

	// Start listening.
	return endpoint.Listen()
}

/*
start as a server listening for incoming requests "127.0.0.1"
*/

func main() {

	//node server
	go server()
	// if err != nil {
	// 	log.Println("Error:", errors.WithStack(err))
	// }

	//webserver to access node state through browser
	// HTTP
	log.Println("start webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "tx_pool_size %d\n%b", len(tx_pool), tx_pool[0])
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

	log.Println("Server running")
}
