FROM alpine
WORKDIR /app
VOLUME ["/data"]
RUN apk add libc6-compat
RUN wget fs.devops.haodai.net/soft/gotty -O /bin/gotty && \
    wget fs.devops.haodai.net/soft/gotty-logs -O /bin/gotty-logs && \
    chmod +x /bin/gotty /bin/gotty-logs

CMD gotty -w gotty-logs 
EXPOSE 8080