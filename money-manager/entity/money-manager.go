package entity

type Fund uint64

type UserId string

type ReserveId string

type ServiceId string

type OrderId string

type Balance struct {
	Current   Fund
	Available Fund
}

type Reserve struct {
	UserId    UserId
	ServiceId ServiceId
	OrderId   OrderId
	Amount    Fund
}
