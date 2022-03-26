package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var (
	ErrDBConnection = errors.New("no DB connection")
	ErrKeyExists    = errors.New("key exists")
)

//type KeyExistsError struct {
//	Key string
//	Err error
//}
//
//func (k *KeyExistsError) Error() string {
//	return fmt.Sprintf("%s: %v", k.Err, k.Key)
//}
//
//func NewKeyExistsError(key string, err error) error {
//	return &KeyExistsError{Key: key, Err: err}
//}

const timeout = 5 * time.Second

type DBRepository struct {
	connection *pgx.Conn
}

func (r *DBRepository) AddBatch(urls []entities.URLItem, userID string) ([]entities.ShortenURL, error) {
	var result []entities.ShortenURL

	ctx := context.Background()
	tx, err := r.connection.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	for _, item := range urls {
		_, err = tx.Exec(ctx, "insert into links (userid, key, url) values ($1, $2, $3)", userID, item.CorrelationID, item.OriginalURL)
		if err != nil {
			return nil, err
		}
		result = append(result, entities.ShortenURL{Key: item.CorrelationID})
	}

	err = tx.Commit(ctx)
	if err != nil {

		return nil, err
	}

	return result, nil
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

	sql := `CREATE TABLE IF NOT EXISTS links(
				userid varchar(16) NOT NULL,
				key varchar(36) NOT NULL, 
				url varchar(450) NOT NULL, 
				PRIMARY KEY (key)
            )`

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

	if err != nil {
		log.Printf("can't insert: [%T] %v", err, err.(*pgconn.PgError).Code)

		if err.(*pgconn.PgError).Code == "23505" {

			log.Printf("Key alrady exists: %s", key)

			return ErrKeyExists
		}

		return err
	}

	log.Printf("insert result: %v", ret)

	return nil
}

func (r *DBRepository) GetAllFor(userID string) []entities.LinkItem {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := `select key, url from links where userid = $1`

	rows, err := r.connection.Query(ctx, sql, userID)
	if err != nil {
		log.Printf("db error: %v", err)
		return nil
	}
	defer rows.Close()

	links := make([]entities.LinkItem, 0)

	for rows.Next() {
		var link entities.LinkItem
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
