FROM harbor.haodai.net/base/alpine:3.7cgo
WORKDIR /app
VOLUME ["/data"]

RUN wget fs.devops.haodai.net/soft/gotty -O /bin/gotty && \
    wget fs.devops.haodai.net/soft/gotty-logs -O /bin/gotty-logs && \
    chmod +x /bin/gotty /bin/gotty-logs

CMD gotty -w gotty-logs 
EXPOSE 8080