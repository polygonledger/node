# Polygon

Polygon is a new blockchain. It is written in golang and uses a novel networking stack built from two primitives: channels and extensible data notation.
On this lower layer an economic system is built - the proof-of-asset algorithm. The lower layer is created such that in principle any sensible economic incentive and consensus algorithm
can be plugged in and more generically a distributed system can be built. Polygon allows for arbitrary message encoding and signing, using new primitives for communication between nodes. This makes it more general as a transaction and communcation platform. See [docs/wp.md](whitepaper for details)

## runing a node

install golang and git

run node:
```./run.sh```

the node will run a peer on port 8888 and a webserver at the same time

with browser go to http://localhost:8080

## client functions

create keys

```cd client && go run client.go -option=createkeys```

 verify signature
 
 ```go run client.go -option=verify```

## telnet 

```
telnet localhost 8888
Connected to localhost.
Escape character is '^]'.
{:REQ STATUS}
{:REP STATUS :data {:Blockheight 1 :LastBlocktime #inst"2020-05-08T18:17:08.881488+07:00" :Servertime #inst"2020-05-08T18:18:04.507472+07:00" :Starttime #inst"2020-05-08T18:17:08.919498+07:00" :Timebehind 55 :Uptime 55}}
```

## testing

```go test```

## contributions

contributions, such as pull requests, bug reports and comments are very welcome

Discord:
https://discord.gg/wf5Qu72

Telegram:
https://t.me/joinchat/Dzif7R1cHnAzulflui53fA

License: MIT license
