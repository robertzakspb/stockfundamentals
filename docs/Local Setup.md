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
EndPoint: invest-public-api.tbank.ru:443
AppName: invest-api-go-sdk
DisableResourceExhaustedRetry: false
DisableAllRetry: false
MaxRetries: 3
```

2. Create a dev.env file and add the following configurations to it:

For macOS:

```text
DB_CONNECTION_STRING=grpc://localhost:2136/local
```

For Linux:

```text
DB_CONNECTION_STRING=grpc://localhost:2136/Root/test
```

3. Add the desired logging mode to the environment file:

```text
LOG_MODE={value}
```

The value may be set to either **CONSOLE** (logs are thus written using the fmt.Println method) or **FILE**. In case the mode is set to **FILE**, the environment file must feature an additional entry containing the **FILE_LOCATION** parameter set to the directory where logs are to be stored.

## Starting the Project Locally

That's it, you are officially good to go. The tasks.json file contains the automatic tasks that start the database on both macOS and Linux and even create backups following the debug session.
