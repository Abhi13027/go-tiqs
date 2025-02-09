package main

import (
	"fmt"
	"os"

	"github.com/Abhi13027/go-tiqs/tiqs"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	userID := os.Getenv("USER_ID")
	password := os.Getenv("PASSWORD")
	totp_key := os.Getenv("TOTP_KEY")
	appID := os.Getenv("APP_ID")
	appSecret := os.Getenv("APP_SECRET")

	client := tiqs.NewClient(appID, appSecret)

	err := client.AutoLogin(userID, password, totp_key)

	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.GetLimits()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

	basketMargin, err := client.GetBasketMargin(tiqs.BasketMarginRequest{
		{
			Exchange:        "NSE",
			Symbol:          "SBIN",
			Quantity:        "100",
			Product:         "I",
			Price:           "0",
			TransactionType: "B",
			OrderType:       "MKT",
		},
		{
			Exchange:        "NSE",
			Symbol:          "REDINGTON",
			Quantity:        "100",
			Product:         "I",
			Price:           "0",
			TransactionType: "B",
			OrderType:       "MKT",
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(basketMargin)

}
