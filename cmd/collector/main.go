package main

import (
	"database/sql"
	"log"
	"time"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=gpu-telemetry-postgresql user=postgres password=postgres dbname=telemetry sslmode=disable")
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	log.Println("Waiting for database...")
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			var exists bool
			err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'telemetry')").Scan(&exists)
			if err == nil && exists {
				log.Println("Database and table ready")
				break
			}
		}
		time.Sleep(2 * time.Second)
	}
	defer db.Close()
	log.Println("collector starting")
	select {}
}
