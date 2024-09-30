package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	storage_suite "gitlab.ozon.dev/chppppr/homework/tests/suite"
)

func TestStorageDBSuite(t *testing.T) {
	suite.Run(t, &storage_suite.StorageDBSuite{})
}
