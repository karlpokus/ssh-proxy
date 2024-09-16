# ssh-proxy
ssh-proxy with idP support without exposing sshd.

# local dev env
The guilty parties: ssh client <-> proxy command <-> proxy server <-> sshd target

requirements:
- docker
- go

````sh
# build all the things
$ make build
# run sshd
$ docker run -d -p 2222:22 sshd
# run proxy-server
$ go run proxy/server/server.go
# run ssh client
$ ssh root@localhost -p 2222 -o ProxyCommand="./bin/cmd %h %p"
````

# todos
- [x] don't expose sshd
- [ ] put sshd_config in docker volume
- [x] make proxy command log to file
- [ ] ssh session timeout
- [ ] proxy command idP login dance
- [ ] mitm
- [ ] set sshd hostname
- [x] re-create in AWS alb? no
- [ ] mv backend into a private network
- [ ] proxy other protocols
- [ ] is http.Mux blocking?
