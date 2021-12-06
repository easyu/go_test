package pojo

import "database/sql/driver"

type User struct {
	Id   int
	Age  int
	Name string `db:"name"`
}

// Value 实现 driver.Valuer 接口
func (u User) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}
