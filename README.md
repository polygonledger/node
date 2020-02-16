# polygon

new blockchain. delegated proof-of-asset algorithm

architecture: accounts, not UTXO. transaction types, scripting through transaction multiplexing

current status: experimental node protocoltxValid(tx)

## run node

node:
```go run node.go```

the node will run a peer and a webserver at the same time

with browser go to http://localhost:8080

## run client

client:
```cd client && go run client.go -option=randomtx```

getbalance:
go run client.go -option=getbalance

## wallet

create keys

```cd client/wallet && go run utils.go -option=createkeys```

 verify signature
 
 ```go run utils.go -option=verify```

## testing

```go test```