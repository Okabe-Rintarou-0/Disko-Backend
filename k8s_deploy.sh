# start prometheus
kubectl create configmap prometheus --from-file=./deployment/documents/prometheus.yml
kubectl apply -f prometheus.yml

# start mysql
kubectl create configmap mysql-init-script --from-file=./deployment/documents/init.sql
kubectl apply -f mysql.yml

# start redis
kubectl apply -f redis.yml

# start disko backend
kubectl apply -f disko.yml
