FROM scratch
ADD gcp-signing-proxy /gcp-signing-proxy
ADD cacert.pem /cacert.pem
CMD ["/gcp-signing-proxy"]
