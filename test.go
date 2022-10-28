package tests

import (
	"os"
	"strings"
)

func PG() string {
	pg := os.Getenv("GOBL_TEST_PG")
	if pg == "" {
		pg = "postgres://localhost:5432"
	}
	return pg + "/gobl_test"
}

func CR() string {
	cr := os.Getenv("GOBL_TEST_CR")
	if cr == "" {
		cr = "postgres://root@localhost:26257"
	}
	return cr + "/gobl_test"
}

func StorageType() string {
	env := strings.ToLower(os.Getenv("GOBL_TEST_STORAGE"))
	switch env {
	case "postgres":
		return "postgres"
	case "cockroach":
		return "cockroach"
	case "sqlite", "":
		return "sqlite"
	default:
		panic("Unknown GOBL_TEST_STORAGE value. Should be one of: pg, cr, sqlite (default)")
	}
}
