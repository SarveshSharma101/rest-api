package datamodels

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	UserId string `json:"userid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Tokens struct {
	Jwt     string
	Refresh string
}

type LoginReq struct {
	UserId string `json:"userId"`
	Pwd    string `json:"password"`
}
