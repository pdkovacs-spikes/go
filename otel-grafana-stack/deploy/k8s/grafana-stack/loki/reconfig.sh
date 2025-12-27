kubectl -n my-grafana delete configmap loki-local-config
kubectl -n my-grafana create configmap loki-local-config --from-file=./local-config.yaml
kubectl -n my-grafana apply -f loki-deployment.yaml
kubectl -n my-grafana scale deployment loki --replicas 0
kubectl -n my-grafana scale deployment loki --replicas 1
