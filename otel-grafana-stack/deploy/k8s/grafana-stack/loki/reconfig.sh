kubectl -n monitoring delete configmap loki-local-config
kubectl -n monitoring create configmap loki-local-config --from-file=./local-config.yaml
kubectl -n monitoring apply -f loki-deployment.yaml
kubectl -n monitoring scale deployment loki --replicas 0
kubectl -n monitoring scale deployment loki --replicas 1
