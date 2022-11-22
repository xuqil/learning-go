package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngin interface {
	// Render 渲染页面
	// tplName 模板的名字，按名索引
	// data 渲染的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateEngin struct {
	T *template.Template
}

func (g *GoTemplateEngin) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}
