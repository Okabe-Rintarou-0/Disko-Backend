# start prometheus
kubectl apply -f ../k8s/prom.yml

# start mysql
kubectl apply -f ../k8s/mysql.yml

# start redis
kubectl apply -f ../k8s/redis.yml

# start disko backend
kubectl apply -f ../k8s/disko.yml
