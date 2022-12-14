package money_manager

import (
	"math/big"
	"regexp"
	"strconv"

	"money-manager/money-manager/entity"
)

func isValidInputMoney(val string, unit string) bool {
	if !isValidMoneyUnit(unit) {
		return false
	}

	if unit == entity.RubUnit {
		return isValidMoneyInRub(val)
	}
	return isValidMoneyInKop(val)
}

func isValidMoneyUnit(str string) bool {
	return str == entity.RubUnit || str == entity.KopUnit
}

func isValidMoneyInKop(str string) bool {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil || val == 0 {
		return false
	}

	return true
}

func isValidMoneyInRub(str string) bool {
	//fund value without kopeyks
	val, err := strconv.ParseUint(str, 10, 64)

	if err != nil {
		checkStr := `^\d+\.?\d{0,2}$`
		regExp := regexp.MustCompile(checkStr)
		if !regExp.MatchString(str) {
			return false
		}

		rubVal, kopVal, ok := splitStrToRubAndKop(str)
		if !ok {
			return false
		}

		bigRub := new(big.Int).SetUint64(rubVal)
		bigRubInKop := bigRub.Mul(big.NewInt(int64(entity.RubValue)), bigRub)
		bigKop := new(big.Int).SetUint64(kopVal)
		maxUint := new(big.Int).SetUint64(^uint64(0))

		totalCop := bigKop.Add(bigRubInKop, bigKop)

		return totalCop.Cmp(maxUint) != 1 && totalCop.Int64() != 0
	}

	return val != 0
}

//isValidAmountSum compares (cur.Money + add.Money) to MaxUint64 and returns:
//
//true, if (cur.Money + add.Money) <= MaxUint64
//false, otherwise
func isValidAmountSum(cur entity.MoneyAmount, add entity.MoneyAmount) bool {
	bigUintCur := new(big.Int).SetUint64(uint64(cur))
	bigUintAdd := new(big.Int).SetUint64(uint64(add))
	bigUintMax := new(big.Int).SetUint64(^uint64(0))

	cmpRes := bigUintCur.Add(bigUintCur, bigUintAdd).Cmp(bigUintMax)
	return cmpRes != 1
}

//isValidDebit returns true if curAmount >= toDebit, otherwise false
func isValidDebit(curAmount entity.MoneyAmount, toDebit entity.MoneyAmount) bool {
	return curAmount >= toDebit
}
