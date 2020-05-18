gen:
	docker run -v ${PWD}:/open-dam  openapitools/openapi-generator-cli generate -i /open-dam/api/openapi.yaml -g go-server -o /open-dam -c /open-dam/config.json --git-user-id open-dam --git-repo-id open-dam-api --enable-post-process-file  -t /open-dam/.openapi-generator/templates/go-server
	# GO_POST_PROCESS_FILE="gofmt -w"   

	docker run -v ${PWD}:/open-dam  openapitools/openapi-generator-cli generate -i /open-dam/api/openapi.yaml -g typescript-node -o /open-dam/client/ts/ --git-user-id open-dam --git-repo-id open-dam-api --additional-properties=npmName=@open-dam/open-dam-api,npmVersion=1.0.0

gen-docs:
	# Requrires [Redoc](https://github.com/Redocly/redoc) 
	# `npm install redoc-cli -g --save`
	redoc-cli bundle api/spec.yaml -o docs/index.html --options.disableSearch --options.hideDownloadButton

build:
	go build -o bin/open-dam-api

build-docker:
	docker build -t open-dam-api .

run:
	docker run -p 8080:8080 open-dam-api
