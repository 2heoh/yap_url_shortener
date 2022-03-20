package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

var (
	ErrDBConnection = errors.New("no DB connection")
)

const timeout = 5 * time.Second

type DBRepository struct {
	connection *pgx.Conn
}

func NewDatabaseRepository(dsn string) Repository {
	repo := &DBRepository{
		connection: connect(dsn),
	}

	repo.init()

	return repo
}

func (r *DBRepository) init() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := "CREATE TABLE IF NOT EXISTS links (userid varchar(16) NOT NULL, key varchar(9) NOT NULL, url varchar(450) NOT NULL, PRIMARY KEY (key))"

	_, err := r.connection.Exec(ctx, sql)
	if err != nil {
		log.Printf("init error: %v", err)
	}

}

func (r *DBRepository) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := `select url from links where key = $1`

	var url string
	err := r.connection.QueryRow(ctx, sql, key).Scan(&url)
	if err != nil {
		log.Printf("can't get url: %v", err)
		return "", err
	}

	return url, nil
}

func (r *DBRepository) Add(key string, url string, userID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := "insert into links (userid, key, url) values ($1, $2, $3)"

	ret, err := r.connection.Exec(ctx, sql, userID, key, url)
	log.Printf("insert result: %v", ret)

	if err != nil {
		log.Printf("can't insert: %v", err)
	}

	return err
}

func (r *DBRepository) GetAllFor(userID string) []LinkItem {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := `select key, url from links where userid = $1`

	rows, err := r.connection.Query(ctx, sql, userID)
	if err != nil {
		log.Printf("db error: %v", err)
		return nil
	}
	defer rows.Close()

	links := make([]LinkItem, 0)

	for rows.Next() {
		var link LinkItem
		err = rows.Scan(&link.ShortURL, &link.OriginalURL)
		if err != nil {
			return nil
		}
		links = append(links, link)
	}

	return links
}

func (r *DBRepository) Ping() error {
	if r.connection == nil {
		return ErrDBConnection
	}

	return r.connection.Ping(context.Background())
}

func connect(dsn string) *pgx.Conn {
	connection, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err.Error())

		return nil
	}

	return connection
}
