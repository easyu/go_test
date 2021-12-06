package main

import (
	"database/sql"
	"fmt"
	"github.com/easyu/go_test/pojo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sql.DB
var dbSQLX *sqlx.DB

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
	//err := initDB() // 调用输出化数据库的函数
	//if err != nil {
	//	fmt.Printf("init pojo failed,err:%v\n", err)
	//	return
	//}
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
	//sqlInjectDemo("'xxx'or id=2")
	//transactionDemo()
	err := initSQLXDB()
	if err != nil {
		fmt.Printf("init sqlxdb failed,err:%v\n", err)
		return
	}
	sqlXQueryRowDemo()
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, age from User where id=?"
	var u pojo.User
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := db.QueryRow(sqlStr, 1).Scan(&u.Id, &u.Name, &u.Age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.Id, u.Name, u.Age)
}

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from User where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u pojo.User
		err := rows.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.Id, u.Name, u.Age)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into User(name, age) values (?,?)"
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
	sqlStr := "update User set age=? where id = ?"
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
	sql := "select id,name,age from User where id > ?"
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
		var u pojo.User
		err := rows.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			return
		}
		fmt.Printf("id:%d,name:%s,age:%d\t\n", u.Id, u.Name, u.Age)
	}
}
func prepareInsertDemo() {
	sqlStr := "insert into User(name,age) values(?,?)"
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
	sqlStr := "update User set name=?,age=? where id = ? "
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
	sqlStr := fmt.Sprintf("select * from User where name=%s", name)
	fmt.Printf("sql:%s\n", sqlStr)
	var u pojo.User
	err := db.QueryRow(sqlStr).Scan(&u.Id, &u.Name, &u.Age)
	if err != nil {
		fmt.Printf("sql query failed,err:%v\n", err)
		return
	}
	fmt.Printf("User:%#v\n", u)
}
func transactionDemo() {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update User set name=320 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}
	if affRow1 == 1 {
		fmt.Println("事务提交了")
		tx.Commit()
	} else {
		fmt.Println("事务回滚了")
		tx.Rollback()
	}
}

func initSQLXDB() (err error) {
	dsn := "root:123456@tcp(192.168.232.100:3306)/sql_test?charset=utf8mb4&parseTime=True"
	dbSQLX, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	dbSQLX.SetMaxOpenConns(20)
	dbSQLX.SetMaxIdleConns(10)
	return
}

// 查询单条数据示例
func sqlXQueryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u pojo.User
	err := dbSQLX.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.Id, u.Name, u.Age)
}
