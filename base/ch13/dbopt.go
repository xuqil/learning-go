package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatalln(err)
	}
	// 记得 Close DB
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	CreateTable(db)
	InsertValue(db)
	//DeleteValue(db)
	UpdateValue(db)
	//QueryRow(db)
	//QueryRows(db)
	//Tx(db)
	PrepareStat(db)
}

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}

	// 验证数据库是否可用
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("数据库连接成功")
	return db, err
}

func CreateTable(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	// 除了 SELECT 语句，都是使用 ExecContext
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS user(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	if err != nil {
		log.Fatalf("建表失败: %v", err)
	}
	log.Println("建表成功")
	cancel()
}

func InsertValue(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	// 使用 ? 作为查询的参数的占位符，防止 SQL 注入
	res, err := db.ExecContext(ctx, "INSERT INTO `user`(`id`, `first_name`, `age`, `last_name`) VALUES (?, ?, ?, ?)",
		1, "Tom", 18, "Jerry")
	if err != nil {
		log.Fatalf("插入数据失败：%v", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("获取 受影响行数 失败：%v", err)
	}
	log.Println("受影响行数", affected)
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("获取 最后插入的ID 失败：%v", err)
	}
	log.Println("最后插入的ID", lastId)
	cancel()
}

func DeleteValue(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	res, err := db.ExecContext(ctx, "DELETE FROM `user` WHERE `id` = ?", 1)
	if err != nil {
		log.Fatalf("删除数据失败：%v", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("获取 受影响行数 失败：%v", err)
	}
	log.Println("受影响行数", affected)
	cancel()
}

func UpdateValue(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	res, err := db.ExecContext(ctx, "UPDATE `user` SET first_name = ? WHERE `id` = ?",
		"Smith", 1)
	if err != nil {
		log.Fatalf("更新数据失败：%v", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("获取 受影响行数 失败：%v", err)
	}
	log.Println("受影响行数", affected)
	cancel()
}

func QueryRow(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	// 查询一行数据（预期只有一行）
	row := db.QueryRowContext(ctx,
		"SELECT `id`, `first_name`, `age`, `last_name` FROM `user` WHERE `id` = ?", 1)
	if row.Err() != nil {
		log.Fatalf("查询一行数据失败：%v", row.Err())
	}
	u := User{}
	// 通过 Scan 方法从结果集中获取一行结果
	// 查询不到，会在 Scan 时返回 sql.ErrNoRows
	err := row.Scan(&u.ID, &u.FirstName, &u.Age, &u.LastName)
	if err != nil {
		log.Fatalf("获取结果集失败：%v", err)
	}
	log.Println("结果：", u.String())
	cancel()
}

func QueryRows(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	// 批量查询
	rows, err := db.QueryContext(ctx,
		"SELECT `id`, `first_name`, `age`, `last_name` FROM `user` WHERE `id` = ?", 1)
	if err != nil {
		log.Fatalf("批量查询数据失败：%v", err)
	}
	users := make([]User, 0)
	for rows.Next() { // 标准迭代器设计
		u := User{}
		// 这里没有数据不会返回 sql.ErrNoRows
		// Scan 支持传入的类型:
		//	*string
		//	*[]byte
		//	*int, *int8, *int16, *int32, *int64
		//	*uint, *uint8, *uint16, *uint32, *uint64
		//	*bool
		//	*float32, *float64
		//	*interface{}
		//	*RawBytes
		//	*Rows (cursor value)
		//	any type implementing Scanner (see Scanner docs)
		if err = rows.Scan(&u.ID, &u.FirstName, &u.Age, &u.LastName); err != nil {
			log.Fatalf("获取结果集失败：%v", err)
		}
		users = append(users, u)
		log.Println("结果：", u.String())
	}
	log.Println("最终结果：", users, "长度：", len(users))
	cancel()
}

func Tx(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	// 开始一个事务
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Fatalf("事务开始失败：%v", err)
	}

	// 使用 ? 作为查询的参数的占位符，防止 SQL 注入
	res, err := tx.ExecContext(ctx, "INSERT INTO `user`(`id`, `first_name`, `age`, `last_name`) VALUES (?, ?, ?, ?)",
		1, "Tom", 18, "Jerry")

	if err != nil {
		log.Println("事务中插入数据失败，开始回滚")
		// 回滚
		err = tx.Rollback()
		if err != nil {
			log.Printf("事务回滚失败：%v", err)
		}
		cancel()
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("获取 受影响行数 失败：%v", err)
	}
	log.Println("受影响行数", affected)
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("获取 最后插入的ID 失败：%v", err)
	}
	log.Println("最后插入的ID", lastId)

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("提交事务失败：%v", err)
	}
	cancel()
}

func PrepareStat(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 提前准备好 SQL 语句和占位符
	stmt, err := db.PrepareContext(ctx, "SELECT `id`, `first_name`, `age`, `last_name` FROM `user` WHERE `id`=?")
	if err != nil {
		log.Fatalf("Prepare error: %v", err)
	}
	// 不用 stmt 时需要关闭
	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(stmt)

	// 执行查询语句，id = 1
	rows, err := stmt.QueryContext(ctx, 1)
	if err != nil {
		log.Fatalf("查询失败：%v", err)
	}
	for rows.Next() { // 标准迭代器设计
		u := User{}
		if err = rows.Scan(&u.ID, &u.FirstName, &u.Age, &u.LastName); err != nil {
			log.Fatalf("获取结果集失败：%v", err)
		}
		log.Println("结果：", u.String())
	}
	cancel()
}

type User struct {
	ID        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func (u User) String() string {
	return fmt.Sprintf("ID: %d FirstName: %s Age: %d LastName: %s",
		u.ID, u.FirstName, u.Age, u.LastName.String)
}
