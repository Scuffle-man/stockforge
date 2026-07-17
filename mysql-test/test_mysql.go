package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接字符串：用户名:密码@tcp(地址:端口)/数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/mysql"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ 打开连接失败：", err)
	}
	defer db.Close()

	// 验证连接
	err = db.Ping()
	if err != nil {
		log.Fatal("❌ 连接失败：", err)
	}

	// 获取 MySQL 版本
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal("❌ 查询版本失败：", err)
	}

	fmt.Println("✅ 连接成功，MySQL 版本：", version)
}