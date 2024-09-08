FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=Secretroot
ENV MYSQL_DATABASE=challange
ENV MYSQL_USER=userservice
ENV MYSQL_PASSWORD=*balance*

COPY ./init.sql /docker-entrypoint-initdb.d/

EXPOSE 3306