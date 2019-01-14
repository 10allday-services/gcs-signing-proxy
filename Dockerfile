FROM alpine:3.5
RUN apk add --no-cache ca-certificates && update-ca-certificates

ADD gcp-signing-proxy /gcp-signing-proxy
ADD cacert.pem /cacert.pem
ENV GOOGLE_APPLICATION_CREDENTIALS=/service_account_key.json
WORKDIR /
CMD ["/gcp-signing-proxy"]
