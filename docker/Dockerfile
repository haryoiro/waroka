FROM mysql/mysql-server:8.0.32

COPY init_db.sql /docker-entrypoint-initdb.d/
COPY mysql.cnf /usr/etc/mysql.cnf


ENV MYSQL_ALLOW_EMPTY_PASSWORD=true

# TZの設定
RUN echo 'Asia/Tokyo' > /etc/timezone && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

CMD ["mysqld"]