package tests

import (
	"os"

	"src.goblgobl.com/utils/typed"
)

func PG() string {
	if pg := os.Getenv("GOBL_TEST_PG"); pg != "" {
		return pg
	}
	return "postgres://localhost:5432/gobl_test"
}

func StorageType() string {
	if os.Getenv("GOBL_TEST_STORAGE") == "pg" {
		return "pg"
	}
	return "sqlite"
}

func StorageConfig() typed.Typed {
	if StorageType() == "pg" {
		return typed.Typed{
			"type": "pg",
			"url":  PG(),
		}
	} else {
		return typed.Typed{
			"type": "sqlite",
			"path": ":memory:",
		}
	}
}
