FROM scratch
ENTRYPOINT ["/porkbun-dyndns-client", "run"]
COPY porkbun-dyndns-client /
