package entity

type User struct {
	Username string `json:"username" bson:"_id"`
	Password string `json:"password" bson:"password"` 
}