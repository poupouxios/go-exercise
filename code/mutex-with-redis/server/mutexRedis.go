package main

import (
	"log"
	"time"
	"net"
	"sync"
	"context"
	"strconv"
	"github.com/bsm/redislock"
	"github.com/joho/godotenv"
	pb "github.com/poupouxios/go-exercise/code/mutex-with-redis/exchangeratepb"
	"google.golang.org/grpc"
)

var waitGroup sync.WaitGroup
var rate float64

type RequestedExchangeRate struct {
	redisLock *redislock.Client
	exchangeRate string
	shouldExit bool
}

type server struct {
	pb.UnimplementedExchangeRateServiceServer
}

func (s *server) FetchRate(ctx context.Context, req *pb.CurrencyRequest) (*pb.CurrencyResponse, error) {
	waitGroup.Add(1)
	exchangeRate := req.GetFromCurrency() + req.GetToCurrency()
	log.Printf("Requesting exchange rate for %s", exchangeRate)
	redisClient := getRedisClient()
	redisLock := redislock.New(redisClient)
	reqExchangeRate := RequestedExchangeRate{redisLock,exchangeRate,false}
	go fetchExchangeRate(reqExchangeRate)
	waitGroup.Wait()
	return &pb.CurrencyResponse{Rate: rate}, nil
}


func main(){

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file %v", err)
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExchangeRateServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func fetchExchangeRate(reqExcRate RequestedExchangeRate){
	if reqExcRate.shouldExit {
		waitGroup.Add(1)
	}
	defer waitGroup.Done()
	redisClient := getRedisClient()
	if !reqExcRate.shouldExit {
		retryStrategy := redislock.LimitRetry(redislock.ExponentialBackoff(100*time.Millisecond, 500*time.Millisecond), 200)
		for {
			lock, err := reqExcRate.redisLock.Obtain(context.Background(),"exchange-rate-lock", 10*time.Second, &redislock.Options{
				RetryStrategy: retryStrategy,
			})
			if err != nil {
				log.Printf("Failed to obtain lock %v", err)
				return
			}
			log.Println("Locked")
			defer func() {
				if err := lock.Release(context.Background()); err != nil {
					log.Printf("Failed to release lock %v", err)
				}
			}()
			break;
		}
	}
	log.Printf("Checking exchange rate for %s", reqExcRate.exchangeRate)

	exchangeRate, err := redisClient.Get(context.Background(),"exchangerate:" + reqExcRate.exchangeRate).Result()
	if err != nil {
		log.Printf("Failed to get exchange rate for %s %v", reqExcRate.exchangeRate, err)
		if(reqExcRate.shouldExit){
			log.Printf("Cannot find the exchange rate requested %s", reqExcRate.exchangeRate)
			return
		}
		populateRedisWithExchangeRates()
		time.Sleep(2 * time.Second)
		reqExcRate.shouldExit = true
		fetchExchangeRate(reqExcRate)
		return
	}else{
		log.Printf("Exchange rate for %s is %s", reqExcRate.exchangeRate,exchangeRate)
	}
	rate,err = strconv.ParseFloat(exchangeRate,64)
	if err != nil {
		log.Printf("Failed to parse exchange rate %v", err)
		rate = 0.05
	}
	log.Println("Unlocked")
	return
}