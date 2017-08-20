#!/bin/bash
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o tweet .
