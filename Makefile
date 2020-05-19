.PHONY: build test

gen:
	docker run -v ${PWD}:/open-dam openapitools/openapi-generator-cli generate -i /open-dam/api/openapi.yaml -g go-server -o /open-dam -c /open-dam/.openapi-generator/config.json --git-user-id open-dam --git-repo-id open-dam-api
	# -e GO_POST_PROCESS_FILE="gofmt -w"
	#  -t /open-dam/.openapi-generator/templates/go-server

gen-docs:
	# Requrires [Redoc](https://github.com/Redocly/redoc) 
	# `npm install redoc-cli -g --save`
	redoc-cli bundle api/openapi.yaml -o docs/index.html --options.disableSearch --options.hideDownloadButton

build:
	go build -o ./bin/open-dam-api ./cmd/open-dam-api/

build-docker:
	docker build -t open-dam-api -f build/Dockerfile .

test:
	go test -v ./internal/...

test-coverage:
	if [ ! -d coverage ]; then mkdir coverage; fi
	go test -coverpkg ./internal/... -coverprofile coverage/coverage.out ./... && go tool cover -html=coverage/coverage.out

run:
	docker run -p 8080:8080 open-dam-api
