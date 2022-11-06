package money_manager

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

const (
	RUB string = "RUB"
	KOP        = "KOP"
)

const (
	KopVal int64 = 1
	RubVal       = 100 * KopVal
)

func makeFunds(fundStr string, unitStr string) (entity.Fund, error) {
	var (
		fnd     entity.Fund
		fundVal uint64
		err     error
	)

	unitVal := strings.ToUpper(unitStr)
	if !isValidInputFund(fundStr, unitVal) {
		return fnd, errors.New("err in moneyManager.AddFundsToUser.makeFunds(): Invalid fund")
	}

	if unitVal == RUB {
		fundVal, err = stringRubToKop(fundStr)
	} else {
		fundVal, err = stringToKop(fundStr)
	}

	fnd = entity.Fund(fundVal)
	return fnd, err
}

func balanceToFund(b entity.Balance) entity.Fund {
	//todo: если у человека все в зарезервированных и они вернутся к нему на счет - будет переполнение. Чего делать?
	return b.Current
}

func stringToKop(str string) (uint64, error) {
	if !isValidFundInKop(str) {
		return 0, errors.New("err in moneyManager.AddFundsToUser.makeFunds.stringToKop(): Invalid fund value in KOP")
	}

	val, _ := strconv.ParseUint(str, 10, 64)

	return val, nil
}

func stringRubToKop(str string) (uint64, error) {
	if !isValidFundInRub(str) {
		return 0, errors.New("err in moneyManager.AddFundsToUser.makeFunds.stringRubToKop(): Invalid fund value")
	}

	rub, kop, _ := splitStrToRubAndKop(str)
	totalKop := rub*uint64(RubVal) + kop

	return totalKop, nil
}

func splitStrToRubAndKop(str string) (uint64, uint64, bool) {
	strParts := strings.Split(str, ".")
	rubStr := strParts[0]

	kopStr := ""
	if len(strParts) > 1 && strParts[1] != "" {
		kopStr = strParts[1]
	}
	if len(kopStr) < 2 {
		kopStr = kopStr + "0"
	}

	rubVal, err := strconv.ParseUint(rubStr, 10, 64)
	if err != nil {
		return 0, 0, false
	}
	kopVal, err := strconv.ParseUint(kopStr, 10, 64)
	if err != nil {
		return 0, 0, false
	}

	return rubVal, kopVal, true
}
