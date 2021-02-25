# Blacklight's CCDC Score Server
A scoreserver for teams to use while training for CCDC. Runs on port `8080`.
Note that security was (ironically) _not_ a core design features, since I had to build this in a week. 
My assumption is that I will never run this on an untrusted network or on the public internet, since
it is designed for use in a virtual lab. Call it an accepted risk.

**Therefore**
* you should probably not run this on an untrusted network
* you should probably not run this on the public internet
* you should probably inform your red team this server is out of scope :)

Enjoy!

## Running with Go
Note: requires a newish version of go that supports go modules. I used 1.15. 

```bash
go get
go run server.go
```

## Running with Docker-Compose
Do not run with only docker - there is no Dockerfile for Mongo, since that's dealt with in the `docker-compose.yml` file.


```bash
docker-compose up
```
