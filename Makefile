CREDS ?= $(shell pwd)/service_account_key.json
TAG := local/gcs-signing-proxy:latest

build:
	CGO_ENABLED=0 gox -osarch="linux/amd64" -output="gcs-signing-proxy"
	docker build --no-cache -t ${TAG} .

run:
	docker run -i -t -v ${CREDS}:/service_account_key.json --env-file my.env --rm -p 8000:8000 ${TAG}

clean:
	rm gcs-signing-proxy*
