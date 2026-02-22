package models

// RedisAccess model
type RedisAccess struct {
	AccessUUID string `json:"access_uuid"`
	UserID     string `json:"user_id"`
}

// Token model
type Token struct {
	AccessToken    string `json:"access_token"`
	AccessUUID     string `json:"access_uuid"`
	AccessExpires  int64  `json:"access_expires"`
	RefreshToken   string `json:"refresh_token"`
	RefreshUUID    string `json:"refresh_uuid"`
	RefreshExpires int64  `json:"refresh_expires"`
}
