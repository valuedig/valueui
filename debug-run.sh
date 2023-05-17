#!/bin/bash

set -e

go build -o bin/valueui cmd/valueui/main.go

./bin/valueui -logtostderr true

