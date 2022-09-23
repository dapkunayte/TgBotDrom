package models

import (
	"errors"
)

var ErrNoRecord = errors.New("users: такого пользователя не найдено")

type Snippet struct {
	ID    int
	State string
	Town  string
}
