package repository

import "github.com/pkg/errors"

var ErrNotFound = errors.New("Err in pgMoneyManagerRepo: requested resource is not found")
