FROM harbor.haodai.net/base/alpine:3.7cgo
WORKDIR /app
VOLUME ["/data"]

COPY gotty-pod /bin/

# COPY /home/wen/.kube/config-new ~/.kube/config

RUN wget fs.haodai.net/soft/gotty -O /bin/gotty && \
    wget fs.haodai.net/soft/kubectl -O /bin/kubectl && \
    wget fs.haodai.net/soft/kubectl-debug -O /bin/kubectl-debug && \
    chmod +x /bin/gotty /bin/kubectl /bin/kubectl-debug

CMD gotty --port 8080 -w --permit-arguments gotty-pod -gitlabtoken=MvPVs7Z56gU2k2ADyR6J
EXPOSE 8080