#!/bin/bash

go get golang.org/x/tools/cmd/cover
go test ./... -cover -v -coverprofile=cov.out
