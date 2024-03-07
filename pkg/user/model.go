package user

import (
	"time"

	"github.com/google/uuid"
)
type CreateSessionParams struct {
  ID            string `json:"id" bson:"_id,omitempty"`
	SessionID     uuid.UUID `json:"session_id" bson:"sessionid,omitempty"`
	UserId       uuid.UUID `json:"user_id" bson:"userid,omitempty"`
	RefreshToken []byte    `json:"refresh_token" bson:"refreshtoken"`
	IsBlocked    bool      `json:"is_blocked" bson:"isblocked"`
	ExpiresAt    time.Time `json:"epires_at" bson:"expiresat"`
}