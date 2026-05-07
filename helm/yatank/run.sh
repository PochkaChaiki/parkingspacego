kubectl delete job yandex-tank-test -n load-testing
kubectl delete configmap yandex-tank-config -n load-testing
kubectl create configmap yandex-tank-config --from-file=load.yaml --from-file=ammo.txt -n load-testing
kubectl apply -f deployment.yaml -n load-testing
