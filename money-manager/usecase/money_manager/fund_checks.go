package money_manager

import (
	"math/big"
	"regexp"
	"strconv"

	"money-manager/money-manager/entity"
)

func isValidInputFund(val string, unit string) bool {
	if !isValidFundUnit(unit) {
		return false
	}

	if unit == RUB {
		return isValidFundInRub(val)
	} else {
		return isValidFundInKop(val)
	}
}

func isValidFundUnit(str string) bool {
	return str == RUB || str == KOP
}

func isValidFundInKop(str string) bool {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil || val == 0 {
		return false
	}

	return true
}

func isValidFundInRub(str string) bool {
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
		bigRubInKop := bigRub.Mul(big.NewInt(RubVal), bigRub)
		bigKop := new(big.Int).SetUint64(kopVal)
		maxUint := new(big.Int).SetUint64(^uint64(0))

		totalCop := bigKop.Add(bigRubInKop, bigKop)

		return totalCop.Cmp(maxUint) != 1 && totalCop.Int64() != 0
	}

	return val != 0
}

//isValidFundSum compares (cur.Amount + add.Amount) to MaxUint64 and returns:
//true, if (cur.Amount + add.Amount) <= MaxUint64
func isValidFundSum(cur entity.Fund, add entity.Fund) bool {
	bigUintCur := new(big.Int).SetUint64(uint64(cur))
	bigUintAdd := new(big.Int).SetUint64(uint64(add))
	bigUintMax := new(big.Int).SetUint64(^uint64(0))

	cmpRes := bigUintCur.Add(bigUintCur, bigUintAdd).Cmp(bigUintMax)
	return cmpRes != 1
}
