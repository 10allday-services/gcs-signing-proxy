FROM alpine:3.5
RUN apk add --no-cache ca-certificates && update-ca-certificates

ADD gcs-signing-proxy /gcs-signing-proxy
ADD cacert.pem /cacert.pem
ENV GOOGLE_APPLICATION_CREDENTIALS=/service_account_key.json
WORKDIR /
CMD ["/gcs-signing-proxy"]
