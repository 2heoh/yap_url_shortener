package repositories

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var (
	ErrDBConnection = errors.New("no DB connection")
	ErrKeyExists    = errors.New("key exists")
)

const (
	timeout      = 5 * time.Second
	workersCount = 5
)

type DBRepository struct {
	connection    *pgx.Conn
	deleteChannel chan entities.DeleteCandidate
	locker        sync.Locker
}

func NewDatabaseRepository(dsn string) Repository {
	deleteChannel := make(chan entities.DeleteCandidate)

	repo := &DBRepository{
		connection:    connect(dsn),
		deleteChannel: deleteChannel,
	}

	repo.init()

	for i := 0; i < workersCount; i++ {
		go func() {

			defer func() {
				if x := recover(); x != nil {
					log.Printf("run time panic: %v", x)
				}
			}()

			log.Printf("start worker...")
			for job := range deleteChannel {
				if err := repo.makeDelete(job); err != nil {
					time.Sleep(time.Second * 2)
					log.Printf("retry...")
					deleteChannel <- job
				}
			}
		}()
	}

	return repo
}

func (r *DBRepository) makeDelete(candidate entities.DeleteCandidate) error {
	log.Printf("\\_/   %v \n", candidate)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := "UPDATE  links SET deleted=true WHERE key=$1 and userid=$2"
	ret, err := r.connection.Exec(ctx, sql, candidate.Key, candidate.UserID)

	if err != nil {
		log.Printf("can't update: %v", err)
		return err
	}

	log.Printf("update result: %v", ret)
	return nil
}

func connect(dsn string) *pgx.Conn {
	connection, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err.Error())

		return nil
	}

	return connection
}

func (r *DBRepository) init() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := `CREATE TABLE IF NOT EXISTS links(
				userid varchar(16) NOT NULL,
				key varchar(36) NOT NULL, 
				url varchar(450) NOT NULL, 
				deleted boolean DEFAULT FALSE,
				PRIMARY KEY (key)
            )`

	_, err := r.connection.Exec(ctx, sql)
	if err != nil {
		log.Printf("init error: %v", err)
	}

}

func (r *DBRepository) DeleteBatch(keys []string, userID string) error {

	existsLinks := r.GetAllFor(userID)

	if (len(existsLinks)) == 0 {
		return errors.New("no links found")
	}

	for _, link := range existsLinks {
		for _, key := range keys {
			if link.ShortURL == key {
				log.Printf("send delete: &{%v, %v}\n", key, userID)
				r.deleteChannel <- entities.DeleteCandidate{Key: key, UserID: userID}
			}
		}
	}

	return nil
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

func (r *DBRepository) Get(key string) (*entities.LinkItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := `select url, deleted from links where key = $1`

	var (
		url     string
		deleted bool
	)
	err := r.connection.QueryRow(ctx, sql, key).Scan(&url, &deleted)

	if err != nil {
		log.Printf("can't get url: %v", err)

		return nil, err
	}

	return &entities.LinkItem{ShortURL: key, OriginalURL: url, IsDeleted: deleted}, nil
}

func (r *DBRepository) Add(key string, url string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sql := "insert into links (userid, key, url) values ($1, $2, $3)"
	ret, err := r.connection.Exec(ctx, sql, userID, key, url)

	if err != nil {
		log.Printf("can't insert: [%T] %v", err, err.(*pgconn.PgError).Code)

		if err.(*pgconn.PgError).Code == "23505" {

			log.Printf("Key already exists: %s", key)

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
