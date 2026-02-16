package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func GetDatabaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return "postgres://projek:projek@localhost:5432/projek?sslmode=disable"
}

func GetRedisAddr() string {
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		return addr
	}
	return "localhost:6379"
}

func NewCtrl(t *testing.T) *gomock.Controller {
	return gomock.NewController(t)
}

func AssertNil(t *testing.T, obj interface{}) {
	require.Nil(t, obj)
}

func AssertNotNil(t *testing.T, obj interface{}) {
	require.NotNil(t, obj)
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	require.Equal(t, expected, actual)
}

func AssertNoError(t *testing.T, err error) {
	require.NoError(t, err)
}

func AssertError(t *testing.T, err error) {
	require.Error(t, err)
}
