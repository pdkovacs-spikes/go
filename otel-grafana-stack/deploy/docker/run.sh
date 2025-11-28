#!/bin/bash

# https://github.com/pdkovacs/forked-quickpizza
grafana_stack_network=forked-quickpizza_default

# The conainer name conveniently includes a subnamespace-like "prefix"
docker run --rm --network $grafana_stack_network --name bitkit_otel-grafana-stack-app bitkit/otel-grafana-stack-app

