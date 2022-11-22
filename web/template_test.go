//go:build e2e

package web

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"testing"
)

func TestLoginPage(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &GoTemplateEngin{
		T: tpl,
	}

	h := NewHTTPServer(ServerWithTemplateEngine(engine))
	h.GET("/login", func(ctx *Context) {
		err := ctx.Render("login.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})
	h.Start(":8081")
}
