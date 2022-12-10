package sql_demo

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	db.Ping()
	//	这里就可以使用 DB 了
	//sql.OpenDB()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 除了 SELECT 语句，都是使用 ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	//	完成了建表
	require.NoError(t, err)

	// 使用 ? 作为查询的参数的占位符，防止 SQL 注入
	res, err := db.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (?, ?, ?, ?)",
		1, "Tom", 18, "Jerry")
	require.NoError(t, err)
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响行数", affected)
	lastId, err := res.LastInsertId()
	log.Println(affected)
	log.Println("最后插入的ID", lastId)

	// 查询一行数据（预期只有一行）
	row := db.QueryRowContext(ctx,
		"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
	require.NoError(t, row.Err())
	tm := TestModel{}
	// 主要要用指针
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.NoError(t, err)
	log.Println(tm)

	row = db.QueryRowContext(ctx,
		"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 10)
	require.NoError(t, row.Err())
	// 查询不到，会在 Scan 时返回 sql.ErrNoRows
	tm = TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	log.Println(err) // sql: no rows in result set
	require.Error(t, sql.ErrNoRows, err)

	// 批量查询
	rows, err := db.QueryContext(ctx,
		"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
	require.NoError(t, err)
	users := make([]TestModel, 0)
	for rows.Next() { // 标准迭代器设计
		tm = TestModel{}
		// 这里没有数据不会返回 sql.ErrNoRows
		// Scan 支持传入的类型
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
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
		users = append(users, tm)
	}
	log.Println(users)
	cancel()

}

func TestTx(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	db.Ping()
	//	这里就可以使用 DB 了
	//sql.OpenDB()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 除了 SELECT 语句，都是使用 ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	//	完成了建表
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)

	// 使用 ? 作为查询的参数的占位符，防止 SQL 注入
	res, err := tx.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (?, ?, ?, ?)",
		1, "Tom", 18, "Jerry")

	if err != nil {
		// 回滚
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		cancel()
		return
	}
	require.NoError(t, err)
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响行数", affected)
	lastId, err := res.LastInsertId()
	log.Println(affected)
	log.Println("最后插入的ID", lastId)

	// 提交事务
	err = tx.Commit()

	cancel()

}

func TestPrepareStatement(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	db.Ping()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 除了 SELECT 语句，都是使用 ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	//	完成了建表
	require.NoError(t, err)

	stmt, err := db.PrepareContext(ctx, "SELECT * FROM `test_model` WHERE `id`=?")
	require.NoError(t, err)
	// id = 1
	rows, err := stmt.QueryContext(ctx, 1)
	require.NoError(t, err)
	for rows.Next() { // 标准迭代器设计
		tm := TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
		log.Println(tm)
	}
	cancel()

}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
