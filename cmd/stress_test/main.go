package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gitlab.ozon.dev/chppppr/homework/internal/clients/manager"
	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
	manager_service "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"gitlab.ozon.dev/chppppr/homework/scripts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	timer := time.Now()
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctxWichCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	conn, err := grpc.NewClient(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	mng_service := manager_service.NewManagerServiceClient(conn)
	mng_client := manager.NewManagerServiceClient(mng_service)
	cmd.SetManagerServiceClient(mng_client)
	cmd.SetContext(ctxWichCancel)

	gofakeit.Seed(time.Now().Second())
	wk := workers.NewWorkers(uint(runtime.NumCPU()))
	log.Printf("Started %v workers\n", runtime.NumCPU())

	add_requests_count := cfg.Test.AddResponsesCount
	give_requests_count := cfg.Test.GiveResponsesCount
	refund_requests_count := cfg.Test.RefundResponsesCount
	return_requests_count := cfg.Test.ReturnResponsesCount
	total_requests_count := add_requests_count + give_requests_count
	total_requests_count += refund_requests_count + return_requests_count

	reqs := scripts.GenerateAddRequests(add_requests_count)
	log.Printf("Generated %d add_requests\n", add_requests_count)

	giveRequests, userIDs := scripts.GenerateGiveRequests(reqs, give_requests_count)
	log.Printf("Generated %d give_request\n", give_requests_count)

	refund_request := scripts.GenerateRefundRequests(giveRequests, userIDs, refund_requests_count)
	log.Printf("Generated %d refund_requests\n", refund_requests_count)

	return_request := scripts.GenerateReturnRequests(refund_request, return_requests_count)
	log.Printf("Generated %d return_requests\n", return_requests_count)

	bad_requsts := 0
	go func() {
		count := 0
		for res := range wk.Results {
			if res.Err != nil {
				bad_requsts++
				log.Println(res.Err)
			}
			count++
			if count%5000 == 0 {
				log.Println("total request: ", count)
			}
		}
	}()

	for i := 0; i < add_requests_count; i++ {
		wk.AddTask(&workers.TaskRequest{
			Func: func() error {
				return mng_client.AddOrder(ctxWichCancel, reqs[i])
			},
		})
	}
	log.Printf("Sended %d add_requests\n", add_requests_count)

	for i := 0; i < give_requests_count; i++ {
		wk.AddTask(&workers.TaskRequest{
			Func: func() error {
				return mng_client.GiveOrders(ctxWichCancel, giveRequests[i])
			},
		})
	}
	log.Printf("Sended %d give_request\n", give_requests_count)

	for i := 0; i < refund_requests_count; i++ {
		wk.AddTask(&workers.TaskRequest{
			Func: func() error {
				return mng_client.Refund(ctxWichCancel, refund_request[i])
			},
		})
	}
	log.Printf("Sended %d refund_requests\n", refund_requests_count)

	for i := 0; i < return_requests_count; i++ {
		wk.AddTask(&workers.TaskRequest{
			Func: func() error {
				return mng_client.Return(ctxWichCancel, return_request[i])
			},
		})
	}
	log.Printf("Sended %d return_requests\n", return_requests_count)

	wk.CloseJobs()
	log.Println("Waiting for all responses")
	wk.Wait()

	fmt.Printf("\n\nAdd requests: %v\n", add_requests_count)
	fmt.Printf("Give requests: %v\n", give_requests_count)
	fmt.Printf("Refund requests: %v\n", refund_requests_count)
	fmt.Printf("Return requests: %v\n", return_requests_count)
	fmt.Printf("Total Sended requests: %v\n", total_requests_count)
	fmt.Printf("Total Bad response: %v\n", bad_requsts)
	fmt.Printf("Total time: %v sec.\n", time.Since(timer).Seconds())
}
