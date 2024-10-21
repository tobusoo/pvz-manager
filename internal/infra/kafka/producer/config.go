package producer

import (
	"github.com/IBM/sarama"
)

func PrepareConfig(opts ...Option) *sarama.Config {
	c := sarama.NewConfig()

	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Retry.Max = 5
	c.Net.MaxOpenRequests = 1
	c.Producer.CompressionLevel = sarama.CompressionLevelDefault
	c.Producer.Compression = sarama.CompressionGZIP
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true

	for _, opt := range opts {
		_ = opt.Apply(c)
	}

	return c
}
