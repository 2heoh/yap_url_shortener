package repositories

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Row struct {
	key    string
	url    string
	userID string
}

type FileURLRepository struct {
	io   *bufio.ReadWriter
	file *os.File
}

func (repo *FileURLRepository) Add(key string, url string, userID string) error {
	items, _ := repo.findRowBy(userID)

	for _, item := range items {
		if item.ShortURL == key {
			return nil
		}
	}

	fmt.Printf("new key: '%s' for '%s'\n", key, url)

	return repo.writeRow(&Row{
		key:    key,
		url:    url,
		userID: userID,
	})
}

func (repo *FileURLRepository) GetAllFor(userID string) []LinkItem {
	rows, err := repo.findRowBy(userID)
	if err != nil {
		return nil
	}

	return rows
}

func NewFileURLRepository(filename string) Repository {

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Producer error : %s", err)
		return nil
	}

	return &FileURLRepository{
		file: file,
		io:   bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file)),
	}
}

func (repo *FileURLRepository) Get(key string) (string, error) {
	row, err := repo.findRowByKey(key)
	if err != nil {
		return "", err
	}

	return row.url, nil
}

func (repo *FileURLRepository) findRowByKey(key string) (*Row, error) {
	_, err := repo.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	for {
		line, err := repo.io.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("id is not found: %s", key)
			}
			return nil, err
		}

		if line != "" {
			row := splitLine(strings.TrimSpace(line))

			if key == row.key {
				return row, nil
			}
		}
	}
}

func (repo *FileURLRepository) findRowBy(userID string) ([]LinkItem, error) {
	_, err := repo.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var items []LinkItem

	for {
		line, err := repo.io.ReadString('\n')
		log.Printf("line: %s", line)

		if err != nil {

			if err == io.EOF {
				return items, nil
			}

			return nil, err
		}

		if line != "" {
			row := splitLine(strings.TrimSpace(line))

			if userID == row.userID {
				items = append(items, LinkItem{ShortURL: row.key, OriginalURL: row.url})
			}
		}
	}
}

func (repo *FileURLRepository) writeRow(r *Row) error {
	line := fmt.Sprintf("%s;%s;%s\n", r.userID, r.key, r.url)

	_, err := repo.io.WriteString(line)
	if err != nil {
		return err
	}

	return repo.io.Flush()
}

func splitLine(line string) *Row {
	parts := strings.Split(line, ";")

	return &Row{
		userID: parts[0],
		key:    parts[1],
		url:    parts[2],
	}
}
