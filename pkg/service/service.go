package service

import (
	"time"

	"github.com/google/uuid"
)
type Autorization interface {
	//CreateUser(user todo.User) (uuid.UUID, error)
	RefreshToken( userUUID uuid.UUID,sessionUUID uuid.UUID, TimeLive time.Duration)(string)
	 GenerateToken(/*username string, password string,*/userUUID uuid.UUID, TimeLive time.Duration, RefTimeLive time.Duration) (string,string)
  ParseToken(token string) (uuid.UUID,uuid.UUID, error)
}

type Service struct {
	Autorization
}
