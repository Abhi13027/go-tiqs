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
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Login successful!")

	// Get user details
	user, err := client.GetUserDetails()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("User Details:", user)

	// Get Quotes for an instrument
	quotes, err := client.GetMarketQuote(3045, "ltp")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Quotes:", quotes)

	// Get margin details for a single order

	order, err := client.GetMargin(
		tiqs.MarginRequest{
			Exchange:        "NSE",
			Token:           "3045",
			Quantity:        "10",
			Price:           "2000",
			OrderType:       "LMT",
			Product:         "I",
			TransactionType: "B",
		})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Order Margin:", order)

	holidays, err := client.GetHolidays()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Holidays:", holidays)

	indexList, err := client.GetIndexList()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Index List:", indexList)

	optionChainSymbol, err := client.GetOptionChainSymbol()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Option Chain Symbol:", optionChainSymbol)
}
