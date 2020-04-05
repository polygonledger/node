# polygon

Polygon is a new blockchain. It uses a delegated proof-of-asset algorithm and is written in golang.

The basic architecture is the use of accounts, not UTXO. Transactions are typed. The scripting occurs through transaction multiplexing, which means there are several modes of transactions. Currently transactions exist as simple structures
which can be imported/exported from JSON.

current status: experimental

## runing a node

install golang and git

node:
```go run node.go```

the node will run a peer on port 8888 and a webserver at the same time

with browser go to http://localhost:8080

## client functions

create keys

```cd client && go run client.go -option=createkeys```

 verify signature
 
 ```go run client.go -option=verify```


## testing

```go test```

## contributions

contributions, such as pull requests, bug reports and comments are very welcome

https://discord.gg/wf5Qu72

License: MIT license
