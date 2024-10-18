package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/chppppr/homework/internal/infra/kafka/consumer_group"
)

func HandleSignals(ctx context.Context, wg *sync.WaitGroup) context.Context {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	sigCtx, cancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() {
		defer signal.Stop(sigterm)
		defer wg.Done()
		defer cancel()

		for {
			select {
			case sig, ok := <-sigterm:
				if !ok {
					log.Printf("[HandleSignals] signal chan closed: %s\n", sig.String())
					return
				}

				log.Printf("[HandleSignals] signal recv: %s\n", sig.String())
				return
			case _, ok := <-sigCtx.Done():
				if !ok {
					fmt.Println("[HandleSignals] context closed")
					return
				}

				log.Printf("[HandleSignals] ctx done: %s\n", ctx.Err().Error())
				return
			}
		}
	}()

	return sigCtx
}

func HandleConsumerGroupErr(ctx context.Context, cg sarama.ConsumerGroup, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case chErr, ok := <-cg.Errors():
				if !ok {
					log.Println("[HandleConsumerGroupErr] error: chan closed")
					return
				}

				log.Printf("[HandleConsumerGroupErr] error: %s\n", chErr)
			case <-ctx.Done():
				log.Printf("[HandleConsumerGroupErr] ctx closed: %s\n", ctx.Err().Error())
				return
			}
		}
	}()
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	ctx := HandleSignals(context.Background(), wg)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := consumer_group.NewConsumerGroupHandler(logger)
	cg, err := consumer_group.NewConsumerGroup(
		cfg.Kafka.Config.Brokers,
		cfg.Kafka.GroupID,
		cfg.Kafka.Topics,
		handler,
		consumer_group.WithOffsetsInitial(sarama.OffsetOldest),
	)
	if err != nil {
		log.Fatal("consumer_group.NewConsumerGroup: ", err)
	}
	defer cg.Close()

	HandleConsumerGroupErr(ctx, cg, wg)
	cg.Run(ctx, wg)
	wg.Wait()
}
