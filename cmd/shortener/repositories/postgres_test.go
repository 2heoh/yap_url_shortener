package repositories_test

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

//

func TestCheckDuplicate(t *testing.T) {
	t.Parallel()

	message := "ERROR: duplicate key value violates unique constraint \"links_pkey\" (SQLSTATE 23505)"

	require.True(t, strings.Contains(message, "23505"))
}

func TestSendDelete(t *testing.T) {

}
