## Creating a local backup of the database
To create a local backup of the database, simply run the following terminal command:
```
ydb -e "grpc://localhost:2136" -d "/local" tools dump -o "/Users/robert/Library/Mobile Documents/com~apple~CloudDocs/Backups/stockfundamentals/$(date +%Y%m%d_%H%M%S)"
```

## Restoring a database from a backup
To restore the local database from a local backup, use the following terminal command:
```
ydb -e "grpc://localhost:2136" -d "/local" tools restore -i "/Users/robert/Library/Mobile Documents/com~apple~CloudDocs/Backups/stockfundamentals/{backup directory name}" -p "."
```