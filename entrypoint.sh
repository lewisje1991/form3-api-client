#!/bin/sh

# We have to delay before running the tests because the API container can't have a health check added to it.
sleep 2;

go test ./... -v -cover