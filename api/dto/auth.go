package dto

import "time"

type Credentials struct {
	Pass  string `form:"pass" json:"pass" validate:"required,min=6"`
	Email string `form:"email" json:"email" validate:"required,email"`
}

type RespToken struct {
	Token       string    `form:"tokenId" json:"tokenId"`
	ExpiresAt   time.Time `form:"tokenTTL" json:"expiresAt"`
	UserID      string    `form:"userId" json:"userId"`
	AccessToken string    `form:"accessToken" json:"accessToken"`
}
