package errors

import "errors"

var ErrUnknownKind = errors.New("unknown kind, should be struct or map")
