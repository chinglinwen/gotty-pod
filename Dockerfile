FROM harbor.haodai.net/base/alpine:3.7cgo
WORKDIR /app
VOLUME ["/data"]

COPY /home/wen/soft/bin/kubectl /bin/
COPY /home/wen/.kube/config-new ~/.kube/config

RUN wget fs.devops.haodai.net/soft/gotty -O /bin/gotty && \
    wget fs.devops.haodai.net/soft/gotty-pod -O /bin/gotty-pod && \
    chmod +x /bin/gotty /bin/gotty-pod

CMD gotty -w gotty-pod 
EXPOSE 8080