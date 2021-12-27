package services

import (
	"fmt"
	"hash/crc32"
)

type Generator interface {
	Generate(string) string
}

type IDGenerator struct{}

func (idgen *IDGenerator) Generate(url string) string {
	crc32q := crc32.MakeTable(0xD5828281)

	return fmt.Sprintf("%08x", crc32.Checksum([]byte(url), crc32q))
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}
