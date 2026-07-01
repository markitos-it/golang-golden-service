package shared

import "errors"

var ErrGoldenNotFound error = errors.New("golden not found")
var ErrGoldenBadRequest error = errors.New("bad request")
var ErrInvalidGoldenId error = errors.New("invalid golden id")
var ErrInvalidGoldenName error = errors.New("invalid golden name")
var ErrInvalidPageNumber error = errors.New("invalid page number")
var ErrInvalidPageSize error = errors.New("invalid page size")
