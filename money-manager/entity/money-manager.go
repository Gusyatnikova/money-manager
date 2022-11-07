package entity

import "time"

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

type ReportPeriod struct {
	Year  int
	Month time.Month
}

type Report []interface{}

type ReportMoneyPerService []ReportMoneyPerServiceRaw

type ReportMoneyPerServiceRaw struct {
	ServiceId ServiceId
	Sum       string
}
