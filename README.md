# Project feed-service

Feed-service is a microservice, containing logic for managing feed-related entities: users and posts.
It provides functionality to fetch and create users, subscribers and posts.
Feed-service uses HTTP API for most of the functionality, however, it also has gRPC support 
for User saving and fetching by username. Also, this service does not implement notification
logic, instead, on required actions, it publishes messages to kafka, which then get processed
by corresponding service.
Some endpoints require client to provide JWT token in authorization header, which then passed to the
[auth-service](https://github.com/Alieksieiev0/auth-service) to read the claims.

## Getting Started

To run this service, just clone it, and start it 
using either [Make Run](#run) or [Make Docker Run](#run-in-docker). 
However, to run properly, it requires a separate microservice with
gRPC connection, read claims. In scope of the [Feed Project](https://github.com/Alieksieiev0/feed-templ)
microservice called [auth-service](https://github.com/Alieksieiev0/auth-service) was used.
Also, to process messages published to the kafka, service called [notification-service](https://github.com/Alieksieiev0/auth-service)
was used

## MakeFile

### Build
```bash
make build
```

### Run
```bash
make run
```

### Run in docker
```bash
make docker-run
```

### Run and rebuild in docker
```bash
make docker-build-n-run
```

### Shutdown docker
```bash
make docker-down
```

### Test
```bash
make test
```

### Clean
```bash
make clean
```

### Proto
```bash
make proto
```

### Live Reload
```bash
make watch
```

## Flags
This application supports startup flags, 
that can be passed to change servers and clients urls. 
However, be careful changing feed-service servers urls 
if you are running it using docker-compose, because by default
only ports 3000 and 4000 are exposed 

### REST server
- Long Name: rest-server
- Default: 3000

### gRPC server
- Name: grpc-server
- Default: 4000

### gRPC client
- Name: grpc-client
- Default: 4001
