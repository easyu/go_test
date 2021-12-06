package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	id   int
	age  int
	name string
}

var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:123456@tcp(192.168.232.100:3306)/sql_test?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := initDB() // 调用输出化数据库的函数
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}
	//queryRowDemo()
	//fmt.Println("-------------------")
	//queryMultiRowDemo()
	//fmt.Println("-------------------")
	//insertRowDemo()
	//fmt.Println("-------------------")
	//updateRowDemo()
	//prepareInsertDemo()
	//prepareQueryDemo()
	//prepareUpdateDemo()
	//prepareQueryDemo()
	sqlInjectDemo("'xxx'or id=2")
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "王五", 38)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 3)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

//预编译Prepare
func prepareQueryDemo() {
	sql := "select id,name,age from user where id > ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Printf(" Prepare failed,err:%v\n", err)
		return
	}
	//预编译资源关闭
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed,err:%v\n", err)
		return
	}
	// rows链接关闭
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			return
		}
		fmt.Printf("id:%d,name:%s,age:%d\t\n", u.id, u.name, u.age)
	}
}
func prepareInsertDemo() {
	sqlStr := "insert into user(name,age) values(?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed,err:%v\n", err)
		return
	}
	defer stmt.Close()
	exec1, err := stmt.Exec("崔斯特", "20")
	if err != nil {
		fmt.Printf("exec failed,err:%v\n", err)
		return
	}
	exec2, err := stmt.Exec("弗雷格斯", "21")
	if err != nil {
		fmt.Printf("exec failed,err:%v\n", err)
		return
	}
	id1, err := exec1.LastInsertId()
	id2, err := exec2.LastInsertId()
	fmt.Printf("insert success! id:[%v,%v]\n", id1, id2)
}
func prepareUpdateDemo() {
	sqlStr := "update user set name=?,age=? where id = ? "
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed,err:%v\n", err)
		return
	}
	defer stmt.Close()
	exec1, err := stmt.Exec("崔斯特", "20", 1)
	if err != nil {
		fmt.Printf("exec failed,err:%v\n", err)
		return
	}
	id1, err := exec1.RowsAffected()
	if id1 > 0 {
		fmt.Println("update success!")
	}
}

// 任何时候我们都不要手动去拼sql
func sqlInjectDemo(name string) {
	sqlStr := fmt.Sprintf("select * from user where name=%s", name)
	fmt.Printf("sql:%s\n", sqlStr)
	var u user
	err := db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("sql query failed,err:%v\n", err)
		return
	}
	fmt.Printf("user:%#v\n", u)
}
