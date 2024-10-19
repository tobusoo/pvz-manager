package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	kafka_suite "gitlab.ozon.dev/chppppr/homework/tests/suite/kafka"
)

func TestSuite(t *testing.T) {
	suite.Run(t, &kafka_suite.KafkaSuite{})
}
