# polygon

polygon is new blockchain written in golang. It uses a delegated proof-of-asset algorithm

The basic architecture is the use of accounts, not UTXO. Transactions are typed. The scripting occurs through transaction multiplexing,
which is means there are several modes of transactions.

current status: experimental

## run node

node:
```go run node.go```

the node will run a peer on port 8888 and a webserver at the same time

with browser go to http://localhost:8080

## run client

client:
```cd client && go run client.go -option=randomtx```

## wallet

create keys

```cd client && go run client.go -option=createkeys```

 verify signature
 
 ```go run client.go -option=verify```

## testing

```go test```