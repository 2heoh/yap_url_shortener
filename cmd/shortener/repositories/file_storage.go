package repositories

import (
	"bufio"
	"errors"
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
	_, err := c.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
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

	return nil, errors.New("id is not found: " + key)
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
		filename: filename,
		reader:   NewConsumer(filename),
		writer:   NewProducer(filename),
	}
}

func NewProducer(filename string) *producer {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Producer error : %s", err)
		return nil
	}

	return &producer{
		file:     file,
		filename: filename,
		writer:   bufio.NewWriter(file),
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
	file     *os.File
	writer   *bufio.Writer
	filename string
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
	filename string
	reader   *consumer
	writer   *producer
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

	if err != nil && err.Error() == "id is not found: "+id {
		fmt.Println("new key, not found")
		return r.writer.WriteRow(&Row{
			key: id,
			url: url,
		})
	}

	return nil
}
