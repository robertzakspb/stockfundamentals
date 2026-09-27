## Problem

YDB returns the following error:

"pool.With failed with 1 attempts: non-retryable error occurred on attempt No.1 (idempotent=false): operation/GENERIC_ERROR (code = 400080, address = localhost:2136, nodeID = 1, issues = [{5:53 => 'extraneous input '/' expecting {<EOF>, ';'}'}]) at `github.com/ydb-platform/ydb-go-sdk/v3/internal/conn.(*grpcClientStream).RecvMsg(grpc_client_stream.go:180)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/query.nextPart(result.go:305)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/query.(*streamResult).nextPart(result.go:280)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/query.newResult(result.go:217)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/query.execute(execute_query.go:158)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/query.(*Session).execute(session.go:184)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/query.(*Session).Query(session.go:291)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/query.do.func1(client.go:337)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/pool.(*Pool).try(pool.go:449)` at `github.com/ydb-platform/ydb-go-sdk/v3/internal/pool.(*Pool).With.func4(pool.go:515)` at `github.com/ydb-platform/ydb-go-sdk/v3/retry.Retry.func1(retry.go:265)` at `github.com/ydb-platform/ydb-go-sdk/v3/retry.opWithRecover(retry.go:432)` at `github.com/ydb-platform/ydb-go-sdk/v3/retry.RetryWithResult(retry.go:370)` at `github.com/ydb-platform/ydb-go-sdk/v3/retry.Retry(retry.go:271)`"

## Solution

The problem stems from the table name not being wrapped with the single quotation marks. Instead of "user/permission", ensure that the
passed table name is "`user/permission/`"
