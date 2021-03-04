FROM prom/busybox:glibc

COPY ./tailscale-exporter /usr/bin/tailscale-exporter

EXPOSE 8080

ENTRYPOINT ["/usr/bin/tailscale-exporter"]
