# polygon

architecture: accounts, not UTXO. transaction types, no scripting

current status: experimental server protocol

## run

server:
```go run server.go```

client:
```cd client && go run client.go```

the client will send transactions to the server, and the server adds the transaction to the tx pool

the server will run a node and a webserver at the same time

with browser go to http://localhost:8080