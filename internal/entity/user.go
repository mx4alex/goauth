package entity

import "time"

type UserDB struct {
	Id       	 string 	`json:"id" bson:"_id"`
	Name   	 	 string 	`json:"name" bson:"name"`
	Username 	 string 	`json:"username" bson:"username"`
	Password 	 string     `json:"password" bson:"password"`
	RefreshToken string     `json:"refresh_token" bson:"refresh_token"`
	ExpiresAt 	 time.Time  `json:"expires_at" bson:"expires_at"`
}

type UserSignUp struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSignIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
