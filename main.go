package main

import (
	"context"
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
	// fmt.Println("Connected to the database.")
	int, _ := os.Hostname()
	db, err := database.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		PGUSER, PGPASSWORD, PGHOST, PGPORT, PGDATABASE, 10), int)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	app := echo.New()
	app.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	app.GET("/users", func(c echo.Context) error {
		query, err := db.Query(ctx, "SELECT * FROM users")
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

		return c.JSON(200, users)
	})

	app.Logger.Fatal(app.Start(":8080"))
	// Use the database connection.
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
