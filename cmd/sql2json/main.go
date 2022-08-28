package sql2json

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/go-sql-driver/mysql"
)

type cli struct {
	DbConnectionUri string `help:"database connection URI, e.g.: 'root@tcp(127.0.0.1:3306)/mydb'" required:"" env:"SQL2JSON_DB_CONNECTION_URI"`

	DbUser     string `help:"user for the database" optional:"" env:"SQL2JSON_DB_USER"`
	DbPassword string `help:"password for the database" optional:"" env:"SQL2JSON_DB_PASSWORD"`
	DbName     string `help:"name of the database" optional:"" env:"SQL2JSON_DB_NAME"`
	Query      string `help:"a sql query fetching the results" required:"" env:"SQL2JSON_QUERY"`
}

func Exec(args []string) {
	// kong expects only actual arguments and not the program itself
	args = args[1:]

	var flags cli

	k, err := kong.New(&flags)
	if err != nil {
		log.Fatalf("cannot parse arguments: %v", err)
	}
	_, err = k.Parse(args)
	if err != nil {
		log.Fatalf("cannot parse arguments: %v", err)
	}

	dsn, err := parseDsnConfig(&flags)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connecting to database")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	encoder := json.NewEncoder(os.Stdout)
	err = execQuery(db, flags.Query, encoder)
	if err != nil {
		log.Fatal(err)
	}
}

func execQuery(db *sql.DB, query string, writer *json.Encoder) error {
	var err error

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("cannot execute query: %w", err)
	}
	defer rows.Close()

	header, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("cannot get result column names: %w", err)
	}

	length := len(header)
	for rows.Next() {
		pointers := make([]interface{}, length)
		row := make([]*string, length)

		for i := range pointers {
			pointers[i] = &row[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return fmt.Errorf("cannot scan row: %w", err)
		}

		if err := writer.Encode(mapRow(header, row)); err != nil {
			return fmt.Errorf("cannot write row: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("cannot scan rows: %w", err)
	}

	return nil
}

func mapRow(headers []string, r []*string) map[string]string {
	mapped := make(map[string]string)

	for i := range r {
		v := r[i]
		if v != nil {
			mapped[headers[i]] = *v
		}
	}
	return mapped
}

func parseDsnConfig(args *cli) (string, error) {
	cfg, err := mysql.ParseDSN(args.DbConnectionUri)
	if err != nil {
		return "", fmt.Errorf("invalid mysql connection URI: %w", err)
	}

	dbName := args.DbName
	if dbName != "" {
		cfg.DBName = dbName
	}

	dbUser := args.DbUser
	if dbUser != "" {
		cfg.User = dbUser
	}

	dbPassword := args.DbPassword
	if dbPassword != "" {
		cfg.Passwd = dbPassword
	}

	return cfg.FormatDSN(), nil
}
