
generate-openapi:
	docker run \
      -v "${PWD}/openapi/spec:/usr/src" \
      -v "${PWD}/openapi/go:/usr/dest" \
      -w /usr/src \
      openapitools/openapi-generator-cli:v5.4.0 generate \
      -i /usr/src/openapi.yaml \
      -g go-gin-server \
      -o /usr/dest \
      --additional-properties=apiPath='',packageName=openapi,enumClassPrefix=true \
      --global-property models

run-locally:
	API_PORT=8008 \
	LOG_LEVEL=info \
	MONGO_URI=mongodb://localhost:21019/?replicaSet=myrs \
	MONGO_DOMAIN_DB=apart_deal_api \
	REDIS_URI=http://localhost:6379 \
	REDIS_DB=2 \
	go run cmd/apart-deal-api/main.go
