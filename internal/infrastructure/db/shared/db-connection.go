package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

var pool = DbConnectionPool{}

type DbConnectionPool struct {
	connections []DbConnection
}

type DbConnection struct {
	db       *ydb.Driver
	occupied bool
}

func ReleaseDriver(driver *ydb.Driver) error {
	for i := range pool.connections {
		if driver == pool.connections[i].db {
			pool.connections[i].occupied = false
			for i := range pool.connections {
				fmt.Println("Connection: ", &pool.connections[i].db, ". Occupied: ", pool.connections[i].occupied)
			}
			return nil
		}
	}

	for i := range pool.connections {
		fmt.Println("Connection: ", &pool.connections[i].db, ". Occupied: ", pool.connections[i].occupied)
	}

	logger.Log("Attempting to release a ydb connection not present in the pool", logger.ERROR)
	return errors.New("Attempting to release a ydb connection not present in the pool")
}

func GetReusableYdbDriver() (*ydb.Driver, error) {
	//Attempting to find an unoccupied connection
	for i := range pool.connections {
		if !pool.connections[i].occupied {
			pool.connections[i].occupied = true
			return pool.connections[i].db, nil
		}
	}

	//If all connections are occupied (or there are none), instantiate a new connection and return it
	driver, err := MakeYdbDriver()
	if err != nil {
		return nil, err
	}
	pool.connections = append(pool.connections, DbConnection{
		db:       driver,
		occupied: true,
	})
	return driver, nil
}

func MakeYdbDriver() (*ydb.Driver, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		logger.Log("Failed to fetch the configuration", logger.ALERT)
		return nil, err
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}

	return db, nil
}
