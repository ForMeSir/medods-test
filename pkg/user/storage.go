package user

import (
	"context"

	"github.com/google/uuid"
)


type Storage interface {
	Create(ctx context.Context, arg CreateSessionParams)  (string, error)
	FindOne(ctx context.Context, id uuid.UUID) (s CreateSessionParams, err error)
  // Update(ctx context.Context, user User) error
   Delete(ctx context.Context, id uuid.UUID) error
	// FindAll(ctx context.Context) (u []User,err error)
}
