package main

import (
	"log"
	"time"
	"net/http"
	"net/url"
	"os"
	"io"
	"context"
	"encoding/json"
)

func populateRedisWithExchangeRates() {
	domain := os.Getenv("BASE_EXCHANGE_URL")
	params := url.Values{}
	params.Add("access_key", os.Getenv("EXCHANGE_API_KEY"))

	finalUrl := domain + "live?" + params.Encode()
	log.Printf("Fetching exchange rate from %s", finalUrl)
	resp, err := http.Get(finalUrl);

	if err != nil {
		log.Printf("Failed to fetch exchange rate %v", err)
		return
	}

	responseData,err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body %v", err)
		return
	}
	//log.Printf("Response body  %v", string(responseData))

	data := map[string]interface{}{}

	err = json.Unmarshal(responseData, &data)
	if err != nil {
		log.Printf("Failed to unmarshal response data %v", err)
		return
	}

	if value, ok := data["quotes"].(map[string]interface{}); ok {
		redisClient := getRedisClient()
		for key, exchangeRate := range value {
			//log.Printf("Exchange rate for %s is %v", key, exchangeRate)
			redisKey := "exchangerate:" + key
			err := redisClient.Set(context.Background(),redisKey, exchangeRate, 6 * time.Second).Err()
			if err != nil {
				log.Printf("Failed to set exchange rate to redis %v", err)
				return
			}
		}
	}
}