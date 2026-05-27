## How to Delete the Redundant ~backup Directory in YDB

### Problem
In case the ```ydb tools dump``` unexpectedly terminates while performing a database backup, YDB will attempt
to create a backup directory titled '~backup_yearMonthDateTimestamp' directly in the database. In case this 
directory needs to be permanently removed from the database, execute the following command:

```
ydb -d "/Root/test" -e "grpc://localhost:2136" scheme rmdir -rf ~backup_20260312T204456
```