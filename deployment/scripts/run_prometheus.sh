docker run -d \
    -p 9090:9090 \
    --name prometheus \
    --restart=always \
    -v /"${PWD}"/../documents/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus
