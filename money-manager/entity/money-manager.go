package entity

type Fund struct {
	Amount uint64
}

type Balance struct {
	Current   Fund
	Available Fund
}

type User struct {
	UserId string
}
