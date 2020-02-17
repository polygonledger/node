# polygon

new blockchain. delegated proof-of-asset algorithm

architecture: accounts, not UTXO. transaction types, scripting through transaction multiplexing

current status: experimental

## run node

node:
```go run node.go```

the node will run a peer on 8888 and a webserver at the same time

with browser go to http://localhost:8080

## run client

client:
```cd client && go run client.go -option=randomtx```

## wallet

create keys

```cd client/wallet && go run utils.go -option=createkeys```

 verify signature
 
 ```go run utils.go -option=verify```

## testing

```go test```