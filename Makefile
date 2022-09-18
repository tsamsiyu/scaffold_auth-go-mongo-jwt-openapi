
generate-openapi:
	find openapi/go -name "*.go" -delete
	docker run \
	  --rm \
      -v "${PWD}/openapi/spec:/usr/src" \
      -v "${PWD}/openapi/go:/usr/dest" \
      -w /usr/src \
      openapitools/openapi-generator-cli:v5.4.0 generate \
      -i /usr/src/openapi.yaml \
      -g go-gin-server \
      -o /usr/dest \
      --additional-properties=apiPath='',packageName=openapi,enumClassPrefix=true \
      --global-property models
	git add openapi/go

run-locally:
	export $(cat .env.local | xargs)
	go run cmd/apart-deal-api/main.go

generate mocks:
	go generate ./...
