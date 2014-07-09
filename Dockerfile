FROM busybox
ADD sidekick /usr/local/bin/sidekick
RUN chmod +x /usr/local/bin/sidekick
ENTRYPOINT ["/usr/local/bin/sidekick"]
