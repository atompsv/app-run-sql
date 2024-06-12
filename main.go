package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jvongxay0308/database-go"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var (
	PGUSER     = os.Getenv("PGUSER")
	PGPASSWORD = os.Getenv("PGPASSWORD")
	PGHOST     = os.Getenv("PGHOST")
	PGPORT     = os.Getenv("PGPORT")
	PGDATABASE = os.Getenv("PGDATABASE")
)

func main() {

	// fmt.Println("Starting the application...")
	// // Connect to the database.
	// db, err := connectTCPSocket()
	// if err != nil {
	// 	log.Fatalf("connectTCPSocket: %v", err)
	// }
	// defer db.Close()

	// fmt.Println("Connected to the database.")
	// dsn := fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable",
	// 	PGHOST, PGUSER, PGPASSWORD, PGDATABASE)
	// db, err := sql.Open("postgres", dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	int, _ := os.Hostname()
	db, err := database.Open("pgx", fmt.Sprintf("postgres://%s:%s@/cloudsql/%s/%s?sslmode=disable&connect_timeout=%d",
		PGUSER, PGPASSWORD, PGHOST, PGDATABASE, 10), int)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// ctx := context.Background()

	query, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("db.Query: %v", err)
	}
	defer query.Close()

	users := make([]User, 0)
	for query.Next() {
		var user User
		if err := query.Scan(&user.ID, &user.Name); err != nil {
			log.Fatalf("query.Scan: %v", err)
		}
		users = append(users, user)
	}

	app := echo.New()
	app.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	app.GET("/users", func(c echo.Context) error {
		return c.JSON(200, users)
	})

	app.Logger.Fatal(app.Start(":8080"))
	// Use the database connection.
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// connectTCPSocket initializes a TCP connection pool for a Cloud SQL
// instance of MySQL.
func connectTCPSocket() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_tcp.go: %s environment variable not set.", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser    = mustGetenv("DB_USER")       // e.g. 'my-db-user'
		dbPwd     = mustGetenv("DB_PASS")       // e.g. 'my-db-password'
		dbName    = mustGetenv("DB_NAME")       // e.g. 'my-database'
		dbPort    = mustGetenv("DB_PORT")       // e.g. '3306'
		dbTCPHost = mustGetenv("INSTANCE_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
	)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("postgres", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?sslmode=disable",
		dbUser, dbPwd, dbTCPHost, dbPort, dbName))
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// ...

	return dbPool, nil
}
