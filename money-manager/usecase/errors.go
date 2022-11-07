package usecase

import "github.com/pkg/errors"

var ErrInvalidUser = errors.New("Invalid user")

var ErrNotFound = errors.New("Requested resource is not found")

var ErrDuplication = errors.New("Resource is already exists")

var ErrInvalidMoney = errors.New("Invalid money")

var ErrMoneyLimitIsExceeded = errors.New("Money limit is exceeded")

var ErrNotEnoughMoney = errors.New("Insufficient funds to complete the operation")

var ErrInvalidReserve = errors.New("Invalid reserve")

var ErrInternalError = errors.New("Internal server error")
