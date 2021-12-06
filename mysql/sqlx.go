package main

import (
	"fmt"
	"github.com/easyu/go_test/pojo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

var sqlxDB *sqlx.DB

func main() {
	initSQLX()
	//sqlXQueryOne()
	//sqlXQueryMore()
	//sqlXInsert()
	//namedExec()
	//namedQuery()
	//sqlXTransaction()
	//u1 := pojo.User{Name: "赵信", Age: 15}
	//u2 := pojo.User{Name: "泰达米尔", Age: 13}
	//
	//users := []pojo.User{
	//	u1, u2,
	//}
	//users2 := []interface{}{u1, u2}
	//batchInsertUsers(users)
	//err := sqlXBatchInsertUsers(users2)
	//if err != nil {
	//	fmt.Printf("err:%v\n", err)
	//	return
	//}
	//defer sqlxDB.Close()
	//batchInsertUserByNamedExec(users)
	/*	ds, err := QueryByIDs([]int{1, 2})
		if err != nil {
			return
		}
		fmt.Printf("%#v\n", ds)
	*/
	data, err := queryAndOrderByIds([]int{1, 2})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", data)
}
func initSQLX() (err error) {
	dsn := "root:123456@tcp(192.168.232.100:3306)/sql_test?charset=utf8mb4&parseTime=True"
	sqlxDB = sqlx.MustConnect("mysql", dsn)
	return
}
func sqlXQueryOne(id int) {
	sqlStr := "select id, name, age from user where id=?"
	var u pojo.User
	err := sqlxDB.Get(&u, sqlStr, id)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.Id, u.Name, u.Age)
}
func sqlXQueryMore() {
	sqlStr := "select * from user where id>?"
	var users []pojo.User
	err := sqlxDB.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("sqlXQueryMore failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}
func sqlXInsert() {
	sqlStr := "insert into  user(name,age) values(?,?)"
	exec, err := sqlxDB.Exec(sqlStr, "蔚奥莱", 16)
	if err != nil {
		fmt.Printf("sqlXInsert failed:%v", err)
		return
	}
	id, err := exec.LastInsertId()
	if err != nil {
		fmt.Printf("get last id failed,err:%v", err)
	}
	fmt.Printf("insert success id:%d", id)
}
func namedExec() (err error) {
	sqlStr := "insert into user(name,age) values(:xm,:nl)"
	exec, err := sqlxDB.NamedExec(sqlStr,
		map[string]interface{}{
			"xm": "jnx",
			"nl": 14,
		},
	)
	if err != nil {
		return err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Printf("insert success id:%d\n", id)
	sqlXQueryOne(int(id))
	return
}

// map 查询
func namedQuery() {
	sqlStr := "select * from user where name=:name"
	rows, err := sqlxDB.NamedQuery(sqlStr, map[string]interface{}{
		"name": "jnx",
	})
	if err != nil {
		fmt.Printf("query failed,err:%v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var user pojo.User
		err := rows.StructScan(&user)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", user)
	}
	// 通过结构体查询
	user := pojo.User{
		Name: "jnx",
	}
	rowsByStruct, err := sqlxDB.NamedQuery(sqlStr, user)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rowsByStruct.Close()
	for rowsByStruct.Next() {
		var u pojo.User
		err := rowsByStruct.StructScan(&u)
		if err != nil {
			fmt.Printf("structScan failed :err%v", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}
func sqlXTransaction() {
	tx, err := sqlxDB.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()
	sqlStr1 := "update  user set age=18 where id=? "
	rs, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		return
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return
	}
	if n != 1 {
		fmt.Printf("update failed ")
		return
	}
}

// 自定义批量插入
func batchInsertUsers(users []pojo.User) error {
	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]interface{}, 0, len(users)*2)
	for _, u := range users {
		valueStrings = append(valueStrings, "(?,?)")
		valueArgs = append(valueArgs, u.Name)
		valueArgs = append(valueArgs, u.Age)
	}
	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO user (name, age) VALUES %s",
		strings.Join(valueStrings, ","))
	fmt.Println(stmt)
	_, err := sqlxDB.Exec(stmt, valueArgs...)
	return err
}

// sqlx 批量插入
// 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func sqlXBatchInsertUsers(users []interface{}) error {
	query, args, _ := sqlx.In(
		"INSERT INTO user (name,age) VALUES (?),(?)",
		users...,
	)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args)  // 查看生成的args
	_, err := sqlxDB.Exec(query, args...)
	return err
}
func batchInsertUserByNamedExec(users []pojo.User) error {
	_, err := sqlxDB.NamedExec("insert into user(name,age) values (:name,:age)", users)
	return err
}

func QueryByIDs(ids []int) (users []pojo.User, err error) {
	// 动态填充id
	query, args, err := sqlx.In("SELECT id, name, age FROM user WHERE id IN (?)", ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	fmt.Println(query)
	query = sqlxDB.Rebind(query)
	fmt.Println(query)
	err = sqlxDB.Select(&users, query, args...)
	return
}
func queryAndOrderByIds(ids []int) (users []pojo.User, err error) {
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	query, args, err := sqlx.In("SELECT id,name, age FROM user WHERE id IN (?) ORDER BY FIND_IN_SET(id, ?)", ids, strings.Join(strIDs, ","))
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = sqlxDB.Rebind(query)
	err = sqlxDB.Select(&users, query, args...)
	return
}
