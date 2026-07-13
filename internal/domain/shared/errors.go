package shared

import (
	"errors"
	"fmt"
)

var ErrGoldenNotFound error = errors.New("golden not found")
var ErrGoldenBadRequest error = errors.New("bad request")
var ErrInvalidGoldenId error = fmt.Errorf("%w: invalid golden id", ErrGoldenBadRequest)
var ErrInvalidGoldenName error = fmt.Errorf("%w: invalid golden name", ErrGoldenBadRequest)
var ErrInvalidGoldenContent error = fmt.Errorf("%w: invalid golden content", ErrGoldenBadRequest)
var ErrInvalidGoldenPoster error = fmt.Errorf("%w: invalid golden poster", ErrGoldenBadRequest)
var ErrInvalidGoldenPosterData error = fmt.Errorf("%w: invalid golden poster data", ErrGoldenBadRequest)
var ErrInvalidPageNumber error = fmt.Errorf("%w: invalid page number", ErrGoldenBadRequest)
var ErrInvalidPageSize error = fmt.Errorf("%w: invalid page size", ErrGoldenBadRequest)
