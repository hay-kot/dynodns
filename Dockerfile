FROM gcr.io/distroless/static
ENTRYPOINT ["/dynodns", "run"]
COPY dynodns /
