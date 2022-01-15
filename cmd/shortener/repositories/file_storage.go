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
	key string
	url string
}

type FileURLRepository struct {
	io   *bufio.ReadWriter
	file *os.File
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

func (repo *FileURLRepository) Add(id string, url string) error {

	_, err := repo.findRowByKey(id)

	if err != nil && err.Error() == "id is not found: "+id {
		fmt.Printf("new key: '%s' for '%s'\n", id, url)
		return repo.writeRow(&Row{
			key: id,
			url: url,
		})
	}

	return nil
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

func (repo *FileURLRepository) writeRow(r *Row) error {
	line := fmt.Sprintf("%s;%s\n", r.key, r.url)

	_, err := repo.io.WriteString(line)
	if err != nil {
		return err
	}

	return repo.io.Flush()
}

func splitLine(line string) *Row {
	parts := strings.Split(line, ";")

	return &Row{
		parts[0],
		parts[1],
	}
}
