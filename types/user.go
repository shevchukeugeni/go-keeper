package types

import (
	"encoding/base64"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"
)

type User struct {
	ID        string    `db:"id"          json:"ID"`
	CreatedAt time.Time `db:"created_at"  json:"created_at"`
	Login     string    `db:"login"       json:"login"`
	Password  string    `db:"password"    json:"password"`
}

func (u *User) ToDB() *User {
	if u == nil {
		return nil
	}

	h := sha3.New512()
	h.Write([]byte(u.Password))

	ret := User{
		ID:       u.ID,
		Login:    u.Login,
		Password: base64.StdEncoding.EncodeToString(h.Sum(nil)),
	}

	if ret.ID == "" {
		ret.ID = uuid.NewV4().String()
	}

	if u.CreatedAt.Unix() <= 0 {
		ret.CreatedAt = time.Now()
	} else {
		ret.CreatedAt = u.CreatedAt
	}

	return &ret
}
