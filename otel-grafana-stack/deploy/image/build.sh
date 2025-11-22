#!/bin/bash

exec_name=otel-grafana-stack
context_dir=deploy/image

cd cmd; go build -o $exec_name .; cd -;

cp cmd/$exec_name $context_dir/
docker build ./$context_dir -t bitkit/$exec_name:latest
