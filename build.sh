export GOPATH=$(pwd)
GOOS=linux GOARCH=amd64 go build sidekick.go

mkdir -p /tmp/sidekick-docker
rsync -a Dockerfile sidekick /tmp/sidekick-docker

docker build -t lokalebasen/sidekick /tmp/sidekick-docker/
