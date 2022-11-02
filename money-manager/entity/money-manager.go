package entity

type User struct {
	UserId string
}

type UserBalance struct {
	Balance uint64
}

type BalanceOperation struct {
	UserId string
	Amount uint64
}
