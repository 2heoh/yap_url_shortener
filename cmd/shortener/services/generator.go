package services

import (
	"fmt"
	"hash/crc32"
)

func GenerateID(url string) string {
	crc32q := crc32.MakeTable(0xD5828281)

	return fmt.Sprintf("%08x", crc32.Checksum([]byte(url), crc32q))
}
