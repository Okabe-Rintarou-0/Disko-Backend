docker run -itd \
    --name mysql \
    --restart=always \
    -e MYSQL_ROOT_PASSWORD=123456 \
    -v /"${PWD}"../documents/init.sql:/docker-entrypoint-initdb.d/init.sql \
    daocloud.io/library/mysql:8