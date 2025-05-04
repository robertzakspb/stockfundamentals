## Setting Up the Project
1. Run the following commands:

```console
go mod tidy
```

## Setting up the configuration files
1. Create a tinkoffAPIconfig.yaml file in the main directory and add the following content to it:

```yaml
AccountId: ""
APIToken: "The token"
EndPoint: invest-public-api.tinkoff.ru:443
AppName: invest-api-go-sdk
DisableResourceExhaustedRetry: false
DisableAllRetry: false
MaxRetries: 3
```
*MAKE SURE TO ADD IT TO GITIGNORE!!*


## Starting the Database Locally
1. Start Docker Desktop
2. Execute the following command to start the container. If the image is not available locally, docker will install it for you:

```console
docker run -d --rm --name ydb-local -h localhost \
  --platform linux/amd64 \
  -p 2135:2135 -p 2136:2136 -p 8765:8765 -p 9092:9092 \
  -v $(pwd)/ydb_certs:/ydb_certs -v $(pwd)/ydb_data:/ydb_data \
  -e GRPC_TLS_PORT=2135 -e GRPC_PORT=2136 -e MON_PORT=8765 \
  -e YDB_KAFKA_PROXY_PORT=9092 \
  ydbplatform/local-ydb:latest
```
4. The database management console will be available at http://localhost:8765/





