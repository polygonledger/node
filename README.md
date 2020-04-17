# polygon

Polygon is a new blockchain. It uses a delegated proof-of-asset algorithm and is written in golang.

Polygon uses accounts to store state. The scripting occurs through transaction multiplexing, which means there are several modes of transactions. Transactions are structures encoded in the [extensible data notation](https://github.com/edn-format/edn).
They are generic code of a certain form which can be understood by the parser.

Polygon allows for arbitrary message encoding and signing, using new primitives for communication between
nodes. This makes it more general as a transaction and communcation platform.

current status: alpha

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
