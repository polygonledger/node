# polygon

architecture: accounts, not UTXO. transaction types, scripting through transaction multiplexing

current status: experimental client-server protocol. once finalized move to peer-to-peer

## run

server:
```go run server.go```

the server will run a node and a webserver at the same time

with browser go to http://localhost:8080

## run client

client:
```cd client && go run client.go -option=randomtx```

the client will send transactions to the server, and the server adds the transaction to the tx pool


getbalance:
go run client.go -option=getbalance

## wallet

create keys

```cd client/wallet && go run utils.go -option=createkeys```

 verify signature
 
 ```go run utils.go -option=verify```

## testing

```go test```