package domain

// User领域对象 是DDD中的聚合根 entity
// BO(business object)
type User struct {
	Email    string
	Password string
}
