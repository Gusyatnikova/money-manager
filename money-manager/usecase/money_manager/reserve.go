package money_manager

import "money-manager/money-manager/entity"

func isValidReserveIds(res entity.Reserve) bool {
	return isValidUserId(res.UserId) && res.OrderId != "" && res.ServiceId != ""
}

func isValidReserveOperation(bal entity.Balance, toReserve entity.Fund) bool {
	isHaveMoney := bal.Available >= toReserve
	curReserve := bal.Current - bal.Available

	return isHaveMoney && isValidFundSum(curReserve, toReserve)
}

func isValidUserId(usr entity.UserId) bool {
	return string(usr) != ""
}
