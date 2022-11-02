package main

import (
	"context"

	"money-manager/money-manager/app"
)

func main() {
	moneyManager := app.NewMoneyManager(context.Background())

	moneyManager.Run()
	moneyManager.ListenForShutdown()
}
