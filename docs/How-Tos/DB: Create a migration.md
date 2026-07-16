## Creating Goose Migrations
In the project, database migrations are handled with the help of a tool called goose. To create a new migration, execute the following command:

```
goose -env="dev.env" create 001_financial_metric_table_initial_fill sql
```