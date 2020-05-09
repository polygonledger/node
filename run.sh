export GIT_COMMIT=$(git rev-list -1 HEAD)
go build -o plynode && ./plynode
