package usecase

import "github.com/pkg/errors"

var ErrInvalidUser = errors.New("Invalid user")

var ErrNotFound = errors.New("requested resources is not found")

var ErrInvalidMoney = errors.New("Invalid money")

var ErrMoneyLimitIsExceeded = errors.New("Money limit is exceeded")

var ErrNotEnoughMoney = errors.New("Insufficient funds to complete the operation")

var ErrInvalidReserve = errors.New("Invalid reserve")
