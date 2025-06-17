package models

import (
	"encoding/json"
	"time"
)

type UserEvent struct {
	UserID    int       `json:"user_id" binding:"required"`
	Username  *string   `json:"username,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Action    string    `json:"action" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

func (ue *UserEvent) EncodeToBinary() []byte {
	data, _ := json.Marshal(ue)
	return data
}

func (ue *UserEvent) DecodeBinary(data []byte) error {
	return json.Unmarshal(data, ue)
}
