package types

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrRecordAlreadyExists = errors.New("record with this key already exists")
