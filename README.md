# Instantlogs

Real time logs aggregation service

## How to run

With golang installed run:

```shell
make run
```

## First steps

1) Start listening new logs...

```shell
curl --no-buffer 'http://localhost:8080/filter?follow'
```

2) In a new shell ingest some logs (any text file will be fine)

```shell
curl http://localhost:8080/ingest --data-binary @/tmp/some.log
```

3) Enjoy

