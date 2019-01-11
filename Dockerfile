FROM scratch
ADD gcp-signing-proxy /gcp-signing-proxy
ADD cacert.pem /cacert.pem
ENV GOOGLE_APPLICATION_CREDENTIALS=/service_account_key.json
WORKDIR /
CMD ["/gcp-signing-proxy"]
