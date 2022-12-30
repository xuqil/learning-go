//go:build v12

package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"leanring-go/orm/internal/errs"
	"leanring-go/orm/internal/valuer"
	"testing"
)

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "invalid columns",
			q:       NewSelector[TestModel](db).Select(C("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "multiple columns",
			q:    NewSelector[TestModel](db).Select(C("FirstName"), C("LastName")),
			wantQuery: &Query{
				SQL: "SELECT `first_name`,`last_name` FROM `test_model`;",
			},
		},
		{
			name: "columns alias",
			q:    NewSelector[TestModel](db).Select(C("FirstName").As("my_name"), C("LastName")),
			wantQuery: &Query{
				SQL: "SELECT `first_name` AS `my_name`,`last_name` FROM `test_model`;",
			},
		},
		{
			name: "avg",
			q:    NewSelector[TestModel](db).Select(Avg("Age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "avg alias",
			q:    NewSelector[TestModel](db).Select(Avg("Age").As("avg_age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`) AS `avg_age` FROM `test_model`;",
			},
		},
		{
			name: "sum",
			q:    NewSelector[TestModel](db).Select(Sum("Age")),
			wantQuery: &Query{
				SQL: "SELECT SUM(`age`) FROM `test_model`;",
			},
		},
		{
			name: "count",
			q:    NewSelector[TestModel](db).Select(Count("Age")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(`age`) FROM `test_model`;",
			},
		},
		{
			name: "max",
			q:    NewSelector[TestModel](db).Select(Max("Age")),
			wantQuery: &Query{
				SQL: "SELECT MAX(`age`) FROM `test_model`;",
			},
		},
		{
			name: "min",
			q:    NewSelector[TestModel](db).Select(Min("Age")),
			wantQuery: &Query{
				SQL: "SELECT MIN(`age`) FROM `test_model`;",
			},
		},
		{
			name:    "aggregate invalid columns",
			q:       NewSelector[TestModel](db).Select(Min("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "multiple aggregate",
			q:    NewSelector[TestModel](db).Select(Min("Age"), Max("Age")),
			wantQuery: &Query{
				SQL: "SELECT MIN(`age`),MAX(`age`) FROM `test_model`;",
			},
		},
		{
			name: "raw expression",
			q:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT `first_name`)")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
			},
		},
		{
			name: "columns in where",
			q:    NewSelector[TestModel](db).Where(C("Id").As("my_id").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = ?;",
				Args: []any{18},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_Build(t *testing.T) {
	db := memoryDB(t)

	testCases := []struct {
		name    string
		builder QueryBuilder

		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "select",
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: NewSelector[TestModel](db).From("`test_model`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from db",
			builder: NewSelector[TestModel](db).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty where",
			builder: NewSelector[TestModel](db).Where(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](db).Where(Not(C("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{18},
			},
		},
		{
			name:    "and",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(18).And(C("FirstName").Eq("Tom"))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "or",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(18).Or(C("FirstName").Eq("Tom"))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) OR (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "invalid column",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(18).Or(C("XXXX").Eq("Tom"))),
			wantErr: errs.NewErrUnknownField("XXXX"),
		},
		{
			name:    "raw expression as predicate",
			builder: NewSelector[TestModel](db).Where(Raw("`id`<?", 18).AsPredicate()),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`id`<?);",
				Args: []any{18},
			},
		},
		{
			name:    "raw expression used in predicate",
			builder: NewSelector[TestModel](db).Where(C("Id").Eq(Raw("`age`+?", 1))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = (`age`+?);",
				Args: []any{1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	// 对应于 query error
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))

	// 对应于 no rows
	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	// data
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Tom", "18", "Jerry")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	//// scan error
	//rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	//rows.AddRow("abc", "Tom", "18", "Jerry")
	//mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	testCases := []struct {
		name string
		s    *Selector[TestModel]

		wantErr error
		wantRes *TestModel
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("XXX").Eq(1)),
			wantErr: errs.NewErrUnknownField("XXX"),
		},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no rows",
			s:       NewSelector[TestModel](db).Where(C("Id").LT(1)),
			wantErr: ErrNoRows,
		},
		{
			name: "data",
			s:    NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantRes: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
		//{
		//	name:    "scan error",
		//	s:       NewSelector[TestModel](db).Where(C("Id").LT(1)),
		//	wantErr: ErrNoRows,
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

type TestModel struct {
	Id int64
	// ""
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func (TestModel) CreateSQL() string {
	return `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`
}

func memoryDB(t *testing.T, opts ...DBOption) *DB {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory",
		// 仅仅用于单元测试
		opts...)
	require.NoError(t, err)
	return db
}

// 在 orm 目录下执行
// go test -bench=BenchmarkQuerier_Get -benchmem -benchtime=10000x
func BenchmarkQuerier_Get(b *testing.B) {
	db, err := Open("sqlite3", fmt.Sprintf("file:benchmark_get.db?cache=shared&mode=memory"))
	if err != nil {
		b.Fatal(err)
	}
	_, err = db.db.Exec(TestModel{}.CreateSQL())
	if err != nil {
		b.Fatal(err)
	}

	res, err := db.db.Exec("INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)"+
		"VALUES (?,?,?,?)", 12, "Deng", 18, "Ming")

	if err != nil {
		b.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		b.Fatal(err)
	}
	if affected == 0 {
		b.Fatal()
	}

	b.Run("unsafe", func(b *testing.B) {
		db.creator = valuer.NewUnsafeValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("reflect", func(b *testing.B) {
		db.creator = valuer.NewReflectValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
