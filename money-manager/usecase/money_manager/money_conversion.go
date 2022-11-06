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

//makeMoney checks if it is possible to convert curAmountStr (curUnitStr) to value in kopeyks
//and returns the resulting value if yes
func makeMoney(curAmountStr string, curUnitStr string) (entity.Fund, error) {
	var (
		fnd       entity.Fund
		curAmount uint64
		err       error
	)

	curUnit := strings.ToUpper(curUnitStr)
	if !isValidInputFund(curAmountStr, curUnit) {
		return fnd, errors.New("err in moneyManager.AddMoneyToUser.makeMoney(): Invalid fund")
	}

	if curUnit == RUB {
		curAmount, err = stringRubToKop(curAmountStr)
	} else {
		curAmount, err = stringToKop(curAmountStr)
	}

	fnd = entity.Fund(curAmount)
	return fnd, err
}

func balanceToFund(b entity.Balance) entity.Fund {
	return b.Current
}

//stringToKop convert kopeyks in string to kopeyks in uint64
func stringToKop(str string) (uint64, error) {
	if !isValidFundInKop(str) {
		return 0, errors.New("err in moneyManager.AddMoneyToUser.makeMoney.stringToKop(): Invalid fund value in KOP")
	}

	val, _ := strconv.ParseUint(str, 10, 64)

	return val, nil
}

//stringRubToKop convert rubles in string to kopeyks in uint64
func stringRubToKop(str string) (uint64, error) {
	if !isValidFundInRub(str) {
		return 0, errors.New("err in moneyManager.AddMoneyToUser.makeMoney.stringRubToKop(): Invalid fund value")
	}

	rub, kop, _ := splitStrToRubAndKop(str)
	totalKop := rub*uint64(RubVal) + kop

	return totalKop, nil
}

//splitStrToRubAndKop extract rubVal, kopVal from string and returns its values and bool flag(true, if no errors was happened)
//
//splitStrToRubAndKop("12.34") returns (12, 34, true)
//splitStrToRubAndKop("12.3") returns (12, 30, true)
func splitStrToRubAndKop(str string) (uint64, uint64, bool) {
	strParts := strings.Split(str, ".")
	rubStr := strParts[0]

	kopStr := ""
	if len(strParts) > 1 && strParts[1] != "" {
		kopStr = strParts[1]
	}
	if len(kopStr) < 2 {
		//12.3 -> 12.30
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
