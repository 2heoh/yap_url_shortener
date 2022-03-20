package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

var (
	ErrDBConnection = errors.New("no DB connection")
)

type DBRepository struct {
	connection *pgx.Conn
}

func (D DBRepository) Get(id string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (D DBRepository) Add(id string, url string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (D DBRepository) GetAllFor(userID string) []LinkItem {
	//TODO implement me
	panic("implement me")
}

func (D DBRepository) Ping() error {
	if D.connection == nil {
		return ErrDBConnection
	}

	return nil
}

func NewDatabaseRepository(dsn string) Repository {
	return &DBRepository{
		connection: connect(dsn),
	}
}

func connect(dsn string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf(fmt.Sprintf("Unable to connect to database: %v", err.Error()))
		return nil
	}
	defer conn.Close(context.Background())

	return conn
}
