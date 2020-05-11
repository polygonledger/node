# Install

the following is for Ubuntu but Mac works almost the same

Install git

Install golang - https://golang.org/

sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update -y
sudo apt install -y golang-go

Check go version

```go version```
go1.13.4 linux/amd64

```cd $HOME && mkdir go && cd go```

set GOPATH with
```export GOPATH=/home/ubuntu/go```

```go get -u github.com/polygonledger/node
cd ~/go/src/github.com/polygonledger/node/
go build
```

Check tests with ```go test```

Check firewall with `sudo ufw status numbered`

sudo ufw enable
sudo ufw allow ssh
sudo apt-get install -y ufw

Run with script
./run.sh