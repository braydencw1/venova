FROM registry.access.redhat.com/ubi10/ubi-minimal:latest

COPY venova-linux-amd64 /usr/local/bin/venova

ENTRYPOINT ["/usr/local/bin/venova"]
