TAG := "local/gcp-signing-proxy:latest"

build:
	CGO_ENABLED=0 go build
	docker build --no-cache -t ${TAG} .

run:
	docker run -i -t --rm -p 8000:8000 ${TAG}

clean:
	rm gcp-signing-proxy*
