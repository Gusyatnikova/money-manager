package usecase

import "github.com/pkg/errors"

var ErrNotFound = errors.New("Requested resource is not found")
var ErrInternalError = errors.New("Internal server error")
var ErrDuplication = errors.New("Resource is already exists")
var ErrInvalidUser = errors.New("Invalid user")
var ErrInvalidMoney = errors.New("Invalid money")
var ErrMoneyLimitIsExceeded = errors.New("Money limit is exceeded")
var ErrNotEnoughMoney = errors.New("Insufficient funds to complete the operation")
var ErrInvalidReserve = errors.New("Invalid reserve")
var ErrInvalidReportInfo = errors.New("Invalid request for report")
