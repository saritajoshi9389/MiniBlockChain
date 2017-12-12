#!/usr/bin/env bash
protoc -I ./ ./db.proto --go_out=plugins=grpc:./go