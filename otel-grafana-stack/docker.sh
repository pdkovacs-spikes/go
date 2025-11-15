#!/bin/bash
cd cmd; go build -o otel-grafana-stack .; cd -;
docker build . -t bitkit/otel-grafana-stack:latest
