package domain

import "time"

// User领域对象 是DDD中的聚合根 entity
// BO(business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Nickname string
	// YYYY-MM-DD
	Birthday time.Time
	AboutMe  string
	Ctime    time.Time
	// 不要组合 万一将来可能有其他同名字段
	WechatInfo WechatInfo
}
