package valuer

import (
	"database/sql"
	"leanring-go/orm/model"
)

type Value interface {
	Field(name string) (any, error)
	SetColumns(rows *sql.Rows) error
}

type Creator func(model *model.Model, entity any) Value
