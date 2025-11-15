#!/bin/bash

# https://github.com/pdkovacs/forked-quickpizza
grafana_stack_network=forked-quickpizza_default

docker run --rm --network $grafana_stack_network --name bitkit bitkit/otel-grafana-stack

