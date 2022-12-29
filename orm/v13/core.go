package orm

import (
	"leanring-go/orm/internal/valuer"
	"leanring-go/orm/model"
)

type core struct {
	model *model.Model

	dialect Dialect
	creator valuer.Creator
	r       model.Registry
	mdls    []Middleware
}
