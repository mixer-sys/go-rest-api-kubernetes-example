git tag -a v0.0.1 -m "attempt one"

git push origin v0.0.1

helm install youtube-stats -f ytapiconfig/myvalues.yaml ./charts/youtube-stats-chart/

minikube kubectl port-forward svc/youtube-stats-youtube-stats-chart 80:80


helm uninstall youtube-stats

minikube addons enable ingress

minikube tunnel &
