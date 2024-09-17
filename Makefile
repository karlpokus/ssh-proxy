init:
	docker network create pnet

build:
	docker build sshd -t sshd
	docker build proxy/server -t proxy-server
	go build -o bin/cmd ./proxy/cmd
