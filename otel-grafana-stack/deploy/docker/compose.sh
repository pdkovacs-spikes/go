#!/bin/bash

grafana_stack_compose_home=${HOME}/github/pdkovacs/forked-quickpizza
dcomp_project=forked-quickpizza

infra_compose_file=$grafana_stack_compose_home/compose.grafana-local-stack.monolithic.yaml
compose_files="-f $infra_compose_file -f deploy/docker/service-with-lb.yaml"
docker compose $compose_files down && docker compose $compose_files up -d
