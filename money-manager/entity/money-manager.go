package entity

import "github.com/oklog/ulid/v2"

type Money struct {
	Value string
	Unit  string
}

type MoneyAmount uint64

type Balance struct {
	Current   MoneyAmount
	Available MoneyAmount
}

type UserId string

type ReserveId string
type ServiceId string
type OrderId string

type Reserve struct {
	UserId      UserId
	ServiceId   ServiceId
	OrderId     OrderId
	MoneyAmount MoneyAmount
}

type Transaction struct {
	id ulid.ULID
}
