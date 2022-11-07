package entity

const (
	RubUnit string = "RUB"
	KopUnit        = "KOP"
)

type MoneyAmount uint64

const (
	KopValue MoneyAmount = 1
	RubValue             = 100 * KopValue
)

type Money struct {
	Value string
	Unit  string
}

type Balance struct {
	Current   MoneyAmount
	Available MoneyAmount
}
