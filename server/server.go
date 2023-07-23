package main

import (
	"database/sql"
	"gomailinglist/jsonapi"
	"gomailinglist/mdb"
	"log"
	"sync"
	"mailinglist/grpcapi"
	"mailinglist/jsonapi"

	"github.com/alexflint/go-arg"
)

var args struct {
	DbPath string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILING_BIND_JSON"`
}

func main(){
	arg.MustParse(&args)
	
	if args.DbPath == "" {
		args.DbPath = "list.db"

	}
	if args.BindJson == "" {
		args.BindJson = ":8080"

	}
	if args.BindGrpc == "" {
		args.BindJson = ":8081"

	}

	log.Printf("Using data base %v \n", args.DbPath)

	db, err := sql.Open("sqlite3", args.DbPath)
	if err != nil {
		log.Fatal(err)

	}
	defer db.Close()

	mdb.TryCreate(db)

	var wg sync.WaitGroup

	wg.Add(1)

	go func(){
		log.Printf("Starting json api server... \n")
		jsonapi.Serve(db, args.BindJson)
		wg.Done()
	}()

	wg.Wait()
}