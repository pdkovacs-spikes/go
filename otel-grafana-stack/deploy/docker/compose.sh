#!/bin/bash

grafana_stack_compose_home=${HOME}/github/pdkovacs/grafana-local-stack
dcomp_project=grafana-local-stack

infra_compose_file=$grafana_stack_compose_home/compose.yaml
compose_files="-f $infra_compose_file -f deploy/docker/service-with-lb.yaml"
docker compose $compose_files down && docker compose $compose_files up -d
