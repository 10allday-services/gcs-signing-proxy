TAG := "local/gcp-signing-proxy:latest"

build:
	CGO_ENABLED=0 gox -osarch="linux/amd64" -output="gcp-signing-proxy"
	docker build --no-cache -t ${TAG} .

run:
	docker run -i -t --rm -p 8000:8000 ${TAG}

clean:
	rm gcp-signing-proxy*
