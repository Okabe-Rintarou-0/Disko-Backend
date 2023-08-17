# stop prometheus
kubectl delete -f ../k8s/prom.yml

# stop mysql
kubectl delete -f ../k8s/mysql.yml

# stop redis
kubectl delete -f ../k8s/redis.yml

# stop disko backend
kubectl delete -f ../k8s/disko.yml
