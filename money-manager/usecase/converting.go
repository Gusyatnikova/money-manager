package usecase

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

func (e *moneyManager) makeFunds(fundStr string, unitStr string) (entity.Fund, error) {
	var (
		fnd     entity.Fund
		fundVal uint64
		err     error
	)

	unitVal := strings.ToUpper(unitStr)
	if !e.isValidFundUnit(unitVal) {
		return fnd, errors.New("err in moneyManager.AddFundsToUser.makeFunds(): Invalid fund unit")
	}

	if unitVal == RUB {
		fundVal, err = e.stringRubToKop(fundStr)
	} else {
		fundVal, err = e.stringToKop(fundStr)
	}

	fnd.Amount = fundVal
	return fnd, err
}

func (e *moneyManager) balanceToFund(b entity.Balance) entity.Fund {
	return entity.Fund{
		Amount: b.Current.Amount,
	}
}

func (e *moneyManager) stringToKop(str string) (uint64, error) {
	if !e.isValidFundInKop(str) {
		return 0, errors.New("err in moneyManager.AddFundsToUser.makeFunds.stringToKop(): Invalid fund value in KOP")
	}

	val, _ := strconv.ParseUint(str, 10, 64)

	return val, nil
}

func (e *moneyManager) stringRubToKop(str string) (uint64, error) {
	if !e.isValidFundInRub(str) {
		return 0, errors.New("err in moneyManager.AddFundsToUser.makeFunds.stringRubToKop(): Invalid fund value")
	}

	rub, kop, _ := e.splitStrToRubAndKop(str)
	totalKop := rub*uint64(RubVal) + kop

	return totalKop, nil
}

func (e *moneyManager) splitStrToRubAndKop(str string) (uint64, uint64, bool) {
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
