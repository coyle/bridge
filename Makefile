.PHONY: build test push service_test tmpl deploy lint clean

BRANCH=$(shell git symbolic-ref --short HEAD)
PROJECT_PATH=/go/src/github.com/coyle/bridge

build:
	docker network create test-net
	docker build -t $(BRANCH)/bridge-server .

run-integration-tests:
	docker run -d \
		--name mongo \
		--network test-net \
		-p 127.0.0.1:27017:27017 \
		mongo

	docker run -d \
		--name=bridge-server \
		--network test-net \
		-e MONGO="mongodb://mongo:27017/bridge" \
		-p 5050:5050 \
		$(BRANCH)/bridge-server

	docker run --rm \
    	-v $(PWD):$(PROJECT_PATH) \
    	-w=$(PROJECT_PATH) \
		--network test-net \
		-e MONGO="mongodb://mongo:27017/bridge" \
    	appleboy/golang-testing \
    	sh -c "coverage testing"

	

clean:
	# cleanup server
	docker stop bridge-server || true
	docker rm bridge-server || true
	# cleanup integration tests
	docker stop bridge-integration-tests || true
	docker rm bridge-integration-tests || true
	# cleanup mongo
	docker stop mongo || true
	docker rm mongo || true
	# cleanup docker network
	docker network rm test-net || true