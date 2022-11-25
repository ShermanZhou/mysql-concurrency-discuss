package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var totalTests int = 20

func main() {
	// run binary in console as $  ./binaryname -p a-password
	var dbPassword = flag.String("p", "password", " provide database password")
	flag.Parse()

	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3310)/dbdev1", *dbPassword))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var personCount int
	err = db.QueryRow("select count(*) from PERSON").Scan(&personCount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Person has count %d, completed sql driver sanity test \n", personCount)

	defaultContext := context.Background()
	var waitGroup sync.WaitGroup

	for i := 0; i < totalTests; i++ {
		go updatePerson1(db, defaultContext, &waitGroup)
		go updatePerson2(db, defaultContext, &waitGroup)
		// make sure both transactions are completed to run next in parallel.
		waitGroup.Wait()
		log.Printf("Parellel update completed %d of %d\n", i+1, totalTests)
	}
	fmt.Printf("All update tests are completed")
}

var inserPhone string = "INSERT into dbdev1.PHONE (PERSON_ID, PHONE) values (?, ?)"
var deletePhone string = "DELETE from dbdev1.PHONE where PERSON_ID = ?"

func updatePerson1(db *sql.DB, ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done() // task will be count done if done or fail
	tx, err := db.BeginTx(ctx, nil)
	log.Printf("[UpdatePerson1] begin tx\n")
	if err != nil {
		log.Printf("[UpdatePerson1] db execution error %t \n", err)
		return
	}
	db.ExecContext(ctx, inserPhone, 1, "416-111-1111")
	db.ExecContext(ctx, deletePhone, 1)
	db.ExecContext(ctx, inserPhone, 1, "416-111-1111")
	db.ExecContext(ctx, deletePhone, 1)
	db.ExecContext(ctx, inserPhone, 1, "416-111-1111")
	db.ExecContext(ctx, deletePhone, 1)
	db.ExecContext(ctx, inserPhone, 1, "416-111-1111")
	db.ExecContext(ctx, deletePhone, 1)
	// increase chance to catch dead lock
	time.Sleep(time.Millisecond * 500)
	err = tx.Commit()
	if err != nil {
		log.Printf("[UpdatePerson1] commit error %t \n", err)
		return
	}
	log.Printf("[UpdatePerson1] done tx\n")
}

func updatePerson2(db *sql.DB, ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done() // task will be count done if done or fail
	tx, err := db.BeginTx(ctx, nil)
	log.Printf("[UpdatePerson2] begin tx\n")
	if err != nil {
		log.Printf("[UpdatePerson2] db execution error %t \n", err)
		return
	}
	db.ExecContext(ctx, inserPhone, 2, "416-222-2222")
	db.ExecContext(ctx, deletePhone, 2)
	db.ExecContext(ctx, inserPhone, 2, "416-222-2222")
	db.ExecContext(ctx, deletePhone, 2)
	db.ExecContext(ctx, inserPhone, 2, "416-222-2222")
	db.ExecContext(ctx, deletePhone, 2)
	db.ExecContext(ctx, inserPhone, 2, "416-222-2222")
	db.ExecContext(ctx, deletePhone, 2)
	// increase chance to catch dead lock
	time.Sleep(time.Millisecond * 500)
	err = tx.Commit()
	if err != nil {
		log.Printf("[UpdatePerson2] commit error %t \n", err)
		return
	}
	log.Printf("[UpdatePerson2] done tx\n")
}
