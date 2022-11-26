//go:build e2e

package web

import (
	"html/template"
	"log"
	"mime/multipart"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &GoTemplateEngin{
		T: tpl,
	}
	h := NewHTTPServer(ServerWithTemplateEngine(engine))
	h.GET("/upload", func(ctx *Context) {
		err := ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})

	fu := FileUploader{
		FileField: "myfile",
		DstPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", "upload", header.Filename)
		},
	}

	h.POST("/upload", fu.Handle())
	h.Start(":8081")
}

func TestDownload(t *testing.T) {
	h := NewHTTPServer()

	fu := FileDownloader{
		Dir: filepath.Join("testdata", "download"),
	}
	h.GET("/download", fu.Handle())
	h.Start(":8081")
}

func TestStaticResourceHandler_Handle(t *testing.T) {
	h := NewHTTPServer()

	s, err := NewStaticResourceHandler(filepath.Join("testdata", "static"))
	if err != nil {
		return
	}
	//localhost:8081/static/xxx.jpg
	h.GET("/static/:file", s.Handle)
	h.Start(":8081")
}
