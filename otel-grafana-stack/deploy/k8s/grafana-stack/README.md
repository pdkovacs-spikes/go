# Install the Grafana stack on Kubernetes (Minikube)

## Grafana

https://grafana.com/docs/grafana/latest/setup-grafana/installation/kubernetes/

## Alloy

https://grafana.com/docs/alloy/latest/set-up/install/kubernetes/
https://grafana.com/docs/alloy/latest/configure/kubernetes/#method-2-create-a-separate-configmap-from-a-file

Change service type and add extra ports:

```
cd alloy
helm upgrade alloy grafana/alloy -n my-grafana -f alloy-values.yaml
```

## Tempo

```
cd tempo
helm install -n my-grafana tempo grafana/tempo
helm upgrade tempo grafana/tempo -n my-grafana -f tempo-values.yaml
```