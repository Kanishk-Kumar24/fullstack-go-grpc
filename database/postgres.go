package database

import (
	"context"
	"database/sql"
	"log"

	"fullstack-go-grpc/database/internals/models"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// NewPostgresDB creates a new Bun DB instance.
// DataBase => SQL(innoba DB), postgress(khud ka engine) ==> innoba DB
// MySQL and Postgree ye dono alag DB hai...
// schema ==> kuch hota hai...
// MySQL => 1 schema => 1 DB => n*tables
// Postgress => 1 schema => n*DB => in each DB => n*tables

// NoSQL => json Data...
// Key value paris ke form mein store hota hai..
// VetorDB => collection => DB
// points => table
// vetor embeedings
// vector cosine simialrity, knn, eculidian distance method...

// query ==> predeifed ==> mapped(query)  treesbinarrysearching => log n => n
// "select * from users" => "select" "*" "from" "users"

func NewPostgresDB(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add a query hook for logging.
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Ping the database to verify the connection.
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database.")
	return db
}

// CreateSchema creates the necessary tables in the database if they don't exist.
func CreateSchema(ctx context.Context, db *bun.DB) error {
	models := []interface{}{
		(*models.User)(nil),
	}

	for _, model := range models {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			return err
		}
	}
	log.Println("Database schema checked and tables are ready.")
	return nil
}
