# ssh-proxy
ssh-proxy with idP support without exposing sshd.

# local dev env
The guilty parties:

````
ssh client
    ^
    |
    v
proxy command
    ^
    |           host network
----v-------------------------
proxy server    private network
    ^
    |
    v
sshd
````

requirements:
- docker
- go

````sh
# run once
$ make init
# build all the things
$ make build
# run sshd
$ docker run -d --name sshd --network pnet sshd
# run proxy-server
$ docker run -p 8080:8080 --network pnet proxy-server
# run ssh client
$ ssh root@sshd -o ProxyCommand="./bin/cmd %h %p"
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
- [x] mv backend into a private network
- [ ] proxy other protocols
- [ ] is http.Mux blocking?
