package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseUrl string
)

func TestMain(m *testing.M) {
	databaseUrl = "user=postgres password=1938 host=localhost dbname=jwt-test-task sslmode=disable"
	os.Exit(m.Run())
}
