package models

// Token model
type Token struct {
	AccessToken    string `bson:"access_token"`
	AccessUUID     string `bson:"access_uuid"`
	AccessExpires  int64  `bson:"access_expires"`
	RefreshToken   string `bson:"refresh_token"`
	RefreshUUID    string `bson:"refresh_uuid"`
	RefreshExpires int64  `bson:"refresh_expires"`
}
