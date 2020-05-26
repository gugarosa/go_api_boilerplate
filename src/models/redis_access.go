package models

// RedisAccess model
type RedisAccess struct {
	AccessUUID string `bson:"access_uuid"`
	UserID     string `bson:"user_id"`
}
