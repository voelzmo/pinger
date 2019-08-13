#!/usr/bin/env bash
go build -o ping-app
./ping-app --config-path examples/config.yml
