FROM scratch
ENTRYPOINT ["/dynodns", "run"]
COPY dynodns /
