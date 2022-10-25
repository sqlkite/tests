package tests

import (
	"os"
	"strings"

	"src.goblgobl.com/utils/typed"
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
	case "pg":
		return "pg"
	case "cr":
		return "cr"
	case "sqlite", "":
		return "sqlite"
	default:
		panic("Unknown GOBL_TEST_STORAGE value. Should be one of: pg, cr, sqlite (default)")
	}
}

func StorageConfig() typed.Typed {
	switch StorageType() {
	case "pg":
		return typed.Typed{
			"type": "pg",
			"url":  PG(),
		}
	case "cr":
		return typed.Typed{
			"type": "cr",
			"url":  CR(),
		}
	case "sqlite":
		return typed.Typed{
			"type": "sqlite",
			"path": ":memory:",
		}
	}
	panic("impossible to reach here")
}
