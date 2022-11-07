package delivery

import "github.com/pkg/errors"

var ErrBadContentType = errors.New("Content-Type application/json is missing")

var ErrBadRequestBody = errors.New("Request body is incorrect")

var ErrBadRequestParams = errors.New("Request parameters are incorrect")
