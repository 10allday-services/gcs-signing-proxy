build:
	go build
	docker build --no-cache -t local/gcp-signing-proxy .

run:
	docker run -i -t --rm -p 8000:8000 local/gcp-signing-proxy:latest
