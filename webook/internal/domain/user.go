package domain

import "time"

// User领域对象 是DDD中的聚合根 entity
// BO(business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}
