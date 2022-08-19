# Instantlogs

Real time logs aggregation service

## How to compile

With golang installed build all:

```shell
make build
```

## How to run the server

Once compiled, just launch the binary:

```shell
./bin/instantlogsd
```

Alternatively, compile and run on the fly:

```shell
make run
```

## How to send logs with the client

### Using std input

```shell
your-binary | ./bin/instantlogs
```

Example, use logs from journalctl:

```shell
journalctl -f | ./bin/instantlogs
```

Or just from a file...

```shell
cat /var/log/*.log | ./bin/instantlogs
```

### Using existing file

```shell
./bin/instantlogs --file /var/log/kern.log
```

Get some sample files from: https://www.secrepo.com/

## First steps with cURL

1) Start listening new logs...

```shell
curl --no-buffer 'http://localhost:8080/filter?follow'
```

2) In a new shell ingest some logs (any text file will be fine)

```shell
curl http://localhost:8080/ingest --data-binary @/tmp/some.log
```

3) Enjoy

## Random notes

Use sample logs from https://github.com/logpai/loghub

