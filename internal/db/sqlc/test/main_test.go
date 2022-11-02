package test

import (
	"database/sql"
	"github.com/hosseintrz/suggestion_api/internal/config"
	db2 "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *db2.Queries
var db *sql.DB
var conf *config.Config

func TestMain(m *testing.M) {
	var err error
	conf, err = config.GetConfig("../../..")
	if err != nil {
		log.Fatalln("error -> ", err.Error())
		return
	}

	db, err = sql.Open(conf.DBConfig.Driver, conf.DBConfig.Source)
	if err != nil {
		log.Fatalln("cannot connect to db -> ", err)
	}

	testQueries = db2.New(db)

	os.Exit(m.Run())

}
