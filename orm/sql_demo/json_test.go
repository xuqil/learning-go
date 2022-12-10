package sql_demo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestJsonColumn_Value(t *testing.T) {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"Name":"Tom"}`), value)
	js = JsonColumn[User]{}
	value, err = js.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)
}

func TestJsonColumn_Scan(t *testing.T) {
	testCases := []struct {
		name    string
		src     any
		wantErr error
		wantVal User
		valid   bool
	}{
		{
			name: "nil",
		},
		{
			name:    "string",
			src:     `{"Name":"Tom"}`,
			wantVal: User{Name: "Tom"},
			valid:   true,
		},
		{
			name:    "bytes",
			src:     []byte(`{"Name":"Tom"}`),
			wantVal: User{Name: "Tom"},
			valid:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			js := &JsonColumn[User]{}
			err := js.Scan(tc.src)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, js.Val)
			assert.Equal(t, tc.valid, js.Valid)
		})
	}
}

func TestJsonColumn_ScanTypes(t *testing.T) {
	jsSlice := JsonColumn[[]string]{}
	err := jsSlice.Scan(`["a", "b", "c"]`)
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, jsSlice.Val)
	val, err := jsSlice.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`["a","b","c"]`), val)

	jsMap := JsonColumn[map[string]string]{}
	err = jsMap.Scan(`{"a":"a value"}`)
	assert.Nil(t, err)
	val, err = jsMap.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"a":"a value"}`), val)
}

type User struct {
	Name string
}

func ExampleJsonColumn_Value() {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(value.([]byte)))
	// Output:
	// {"Name":"Tom"}
}

func ExampleJsonColumn_Scan() {
	js := JsonColumn[User]{}
	err := js.Scan(`{"Name":"Tom"}`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(js.Val)
	// Output:
	// {Tom}
}

type UserJson struct {
	ID   int
	Name string
}

func TestJsonColumn_Crud(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	db.Ping()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 创建一个表，其中 name 字段的类型为 json
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS user_json(
    id INTEGER PRIMARY KEY,
    name JSON
)
`)
	//	完成了建表
	require.NoError(t, err)

	js := JsonColumn[UserJson]{Valid: true, Val: UserJson{ID: 1, Name: "Tom"}}
	// 将 JsonColumn 插入数据库
	res, err := db.ExecContext(ctx, "INSERT INTO `user_json`(`id`, `name`) VALUES (?, ?)",
		js.Val.ID, js)
	require.NoError(t, err)
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响行数", affected)
	lastId, err := res.LastInsertId()
	log.Println(affected)
	log.Println("最后插入的ID", lastId)

	// 查询一行数据（预期只有一行）
	row := db.QueryRowContext(ctx,
		"SELECT `name` FROM `user_json` WHERE `id` = ?", 1)
	require.NoError(t, row.Err())
	js2 := JsonColumn[UserJson]{}
	// 主要要用指针
	var name string
	err = row.Scan(&name)
	require.NoError(t, err)
	err = js2.Scan(name)
	require.NoError(t, err)
	assert.Equal(t, `{"ID":1,"Name":"Tom"}`, name)
	log.Println(name) // {"Name":"Tom"}
	assert.Equal(t, UserJson{ID: 1, Name: "Tom"}, js2.Val)
	log.Println(js2.Val)
	cancel()

}
