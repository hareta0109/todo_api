FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD root
ENV MYSQL_DATABASE task_db
ENV MYSQL_USER user
ENV MYSQL_PASSWORD pass

COPY ./build/db/config/my.conf /etc/mysql/conf.d/my.cnf
