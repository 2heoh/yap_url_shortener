package services

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIDGenerator_GenerateShouldReturnCRC32OfPassedUrl(t *testing.T) {

	generator := NewIDGenerator()

	result := generator.Generate("https://example.com")

	require.Equal(t, "96248650", result)
}
