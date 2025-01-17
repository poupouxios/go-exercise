package main

import (
	"log"
	"context"
	"time"
	"google.golang.org/grpc"
	pb "github.com/poupouxios/go-exercise/code/mutex-with-redis/exchangeratepb"
)

func main(){
	conn, err := grpc.Dial("172.17.0.4:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	defer conn.Close()

	client := pb.NewExchangeRateServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exchangeRatesRequest := []string{"USDEUR","USDGBP","USDGGP"}

	count := 1000

	for i := count; i > 0; i-- {
		for _,exchangeRate := range exchangeRatesRequest {
			toCurrency := exchangeRate[3:]
			log.Printf("Requesting exchange rate from grpc client for %s", exchangeRate)
			req := &pb.CurrencyRequest{FromCurrency: "USD", ToCurrency: toCurrency}
			resp, err := client.FetchRate(ctx, req)
			if err != nil {
				log.Fatalf("Failed to fetch exchange rate %v", err)
			}
			log.Printf("Exchange rate for %s is %v", exchangeRate, resp.GetRate())
		}
	}
}
