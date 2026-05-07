helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace --set graphana.adminPassword=admin123 --set prometheus.prometheusSpec.
serviceMonitorSelectorNilUsesHelmValues=false
