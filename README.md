# Polygon

Polygon is a new blockchain. It is written in golang and uses a novel networking stack built from two primitives: channels and extensible data notation.
On this lower layer an economic system is built - the proof-of-asset algorithm. The lower layer is created such that in principle any sensible economic incentive and consensus algorithm
can be plugged in and more generically a distributed system can be built. Polygon allows for arbitrary message encoding and signing, using new primitives for communication between nodes. This makes it more general as a transaction and communcation platform. See [whitepaper for details](https://github.com/polygonledger/docs/blob/master/whitepaper.md)

## runing a node

install golang and git

run node:
```./run.sh```

see [install docs](https://github.com/polygonledger/docs/blob/master/install.md)
see also [telnet](https://github.com/polygonledger/docs/blob/master/telnet.md)

## client functions

create keys

```cd client && go run client.go -option=createkeys```

 verify signature
 
 ```go run client.go -option=verify```


## testing

```go test ./...```

## contributions

contributions, such as pull requests, bug reports and comments are very welcome

Discord:
https://discord.gg/wf5Qu72

Telegram:
https://t.me/joinchat/Dzif7R1cHnAzulflui53fA

License: MIT license
