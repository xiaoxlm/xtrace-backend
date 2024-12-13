FROM debian:stable-slim

EXPOSE 17000

WORKDIR /app
COPY ./bin/n9e-amd-linux /app/n9e
COPY ./etc /app/etc/
COPY ./docker-entrypoint.sh /app/

RUN chmod u+x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh"]