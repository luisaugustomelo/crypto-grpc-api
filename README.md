# klever

## About
### This is a gGRPC API written in Go Lang. The objective of this project is to provide a C(reat)R(ead)U(pdate)D(elete) of cryptocurrencies and provide routes for these coins to be voted up or down. This application relies on the concept of routes and streams to receive real-time data from your modifications.

## Requirements
- Git
- Docker
- Golang@1.16.5

## Optional

### The application has a migration system, to fill the cryptocurrencies of collection from database klever it's necessary run:

```bash
make migrate state=up #up or down
```
### How to install and run

1. Clone Repository

```bash
https://github.com/luisaugustomelo/crypto
```

2. Run docker-compose with Application and Mongodb

```bash
make go #It's starts and/or install application
```

3. Run tests

```bash
make test
```
