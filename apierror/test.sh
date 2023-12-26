#!/bin/bash
go test go-blockchain-server/apierror -v -count 1 -run "TestHandler"
