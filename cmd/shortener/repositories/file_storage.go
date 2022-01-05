package repositories

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type consumer struct {
	file   *os.File
	reader *bufio.Reader
}

type Row struct {
	key string
	url string
}

func (c *consumer) FindByKey(key string) (*Row, error) {

	lines, err := ioutil.ReadAll(c.file)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(lines), "\n") {
		if line != "" {
			row := splitLine(line)

			if key == row.key {
				return row, nil
			}
		}
	}

	return nil, ErrNotFound
}

func splitLine(line string) *Row {
	parts := strings.Split(line, ";")

	return &Row{
		parts[0],
		parts[1],
	}
}

func NewFileURLRepository(filename string) *FileURLRepository {
	return &FileURLRepository{
		reader: NewConsumer(filename),
		writer: NewProducer(filename),
	}
}

func NewProducer(filename string) *producer {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Producer error : %s", err)
		return nil
	}

	return &producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}
}

func NewConsumer(filename string) *consumer {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Consumer error: %s", err)
		return nil
	}
	return &consumer{
		file:   file,
		reader: bufio.NewReader(file),
	}
}

type producer struct {
	file   *os.File
	writer *bufio.Writer
}

func (p *producer) WriteRow(r *Row) error {
	line := fmt.Sprintf("%s;%s\n", r.key, r.url)
	_, err := p.writer.WriteString(line)

	if err != nil {
		return err
	}

	return p.writer.Flush()
}

type FileURLRepository struct {
	reader *consumer
	writer *producer
}

func (r *FileURLRepository) Get(key string) (string, error) {
	row, err := r.reader.FindByKey(key)
	if err != nil {
		return "", err
	}

	return row.url, nil
}

func (r *FileURLRepository) Add(id string, url string) error {

	_, err := r.reader.FindByKey(id)

	if err == ErrNotFound {

		return r.writer.WriteRow(&Row{
			key: id,
			url: url,
		})
	}

	return nil
}
