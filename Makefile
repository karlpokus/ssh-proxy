build:
	docker build sshd -t sshd
	go build -o bin/cmd ./proxy/cmd
