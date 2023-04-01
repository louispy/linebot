build:
	set GOOS="linux" && set GOARCH="amd64" && go build -o bin/main main.go

deploy-dev:
	serverless deploy --config serverless.dev.yml