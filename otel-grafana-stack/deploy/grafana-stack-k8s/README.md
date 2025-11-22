# Install the Grafana stack on Kubernetes

This an incomplete attempt at having these services in a local minikube cluster.

## Grafana

https://grafana.com/docs/grafana/latest/setup-grafana/installation/kubernetes/

## Alloy

https://grafana.com/docs/alloy/latest/set-up/install/kubernetes/
https://grafana.com/docs/alloy/latest/configure/kubernetes/#method-2-create-a-separate-configmap-from-a-file

Change service type and add extra ports:

```
helm upgrade alloy grafana/alloy -n my-grafana -f alloy-values.yaml
```
