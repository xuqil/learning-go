package template_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"strings"
	"testing"
	"text/template"
)

// 点号
func TestTemplateDot(t *testing.T) {
	// {{.}} 中的英文点号“.”相当于 tmpl.Execute(buf, "Jerry") 的第二个参数 Jerry
	// 因此会渲染成：hello Jerry
	tmpl, err := template.New("test").Parse("hello {{.}}")
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, "Jerry")
	assert.NoError(t, err)
	want := "hello Jerry"
	assert.Equal(t, want, buf.String())

	type Inventory struct {
		Material string
		Count    uint
	}
	sweaters := Inventory{"wool", 17}
	// 跟上述一样，英文点号是 tmpl.Execute(buf, sweaters) 的第二个参数
	// 所以 {{.Count}} 会被替换才 sweaters.Count, {{.Material}} 会被替换成 sweaters.Material
	tmpl, err = template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	require.NoError(t, err)
	buf = &bytes.Buffer{}
	err = tmpl.Execute(buf, sweaters)
	assert.NoError(t, err)
	want = "17 items are made of wool"
	assert.Equal(t, want, buf.String())
}

// 作用域
func TestTemplateScope(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name: "Jerry",
		Emails: []string{
			"jerry@qq.com",
			"jerry@163.com",
		},
	}
	// {{.}} 这个.的作用于在 range 循环里面，因此是 Emails 的内容
	tmpl, err := template.New("test").Parse(`{{.Name}}的邮箱地址：
{{range .Emails}}
{{.}}
{{end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址：

jerry@qq.com

jerry@163.com
`
	assert.Equal(t, want, buf.String())
}

// 去除空白
func TestTemplateSpace(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name: "Jerry",
		Emails: []string{
			"jerry@qq.com",
			"jerry@163.com",
		},
	}
	tmpl, err := template.New("test").Parse(`{{ .Name }}的邮箱地址：
{{- range .Emails}}
{{.}}
{{- end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址：
jerry@qq.com
jerry@163.com`
	assert.Equal(t, want, buf.String())
}

func TestTemplateComment(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name: "Jerry",
		Emails: []string{
			"jerry@qq.com",
			"jerry@163.com",
		},
	}
	tmpl, err := template.New("test").Parse(`{{ .Name }}的邮箱地址：
{{- /* 这个“-”的作用于是去除前面的空白，这一行注释也会占一行空白 */}}
{{/* 我是注释，没有去掉空白行 */}}
{{- range .Emails}}
{{.}}
{{- end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址：

jerry@qq.com
jerry@163.com`
	assert.Equal(t, want, buf.String())
}

func TestTemplatePipeline(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name: "Jerry",
		Emails: []string{
			"jerry@qq.com",
			"jerry@163.com",
		},
	}
	tmpl, err := template.New("test").Parse(`{{ .Name | printf "%s的邮箱地址:" }}
{{- range .Emails}}
{{. | printf "- %s" }}
{{- end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址:
- jerry@qq.com
- jerry@163.com`
	assert.Equal(t, want, buf.String())
}

func TestTemplateVariable(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name: "Jerry",
		Emails: []string{
			"jerry@qq.com",
			"jerry@163.com",
		},
	}
	tmpl, err := template.New("test").Parse(`{{ .Name }}的邮箱地址：
{{range $index, $email := .Emails}}
{{- printf "%d. %s\n" $index $email }}
{{- end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址：
0. jerry@qq.com
1. jerry@163.com
`
	assert.Equal(t, want, buf.String())
}

func TestTemplatePipelineVariable(t *testing.T) {
	tmpl, err := template.New("test").Parse(`{{"\"output\""}}
	A string constant.
{{printf "%q" "output"}}
	A function call.
{{"output" | printf "%q"}}
	A function call whose final argument comes from the previous
	command.
{{printf "%q" (print "out" "put")}}
	A parenthesized argument.
{{"put" | printf "%s%s" "out" | printf "%q"}}
	A more elaborate call.
{{"output" | printf "%s" | printf "%q"}}
	A longer chain.
{{with "output"}}{{printf "%q" .}}{{end}}
	A with action using dot.
{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
	A with action that creates and uses a variable.
{{with $x := "output"}}{{printf "%q" $x}}{{end}}
	A with action that uses the variable in another action.
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
	The same, but pipelined.`)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, nil)
	if err != nil {
		panic(err)
	}
}

func TestTemplateCondition(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name: "Jerry",
		Emails: []string{
			"jerry@qq.com",
			"jerry@163.com",
		},
	}
	tmpl, err := template.New("test").Parse(`{{ .Name }}的邮箱地址:
{{range $index, $email := .Emails}}
{{- if eq $index 0 -}}
第一个 email:{{$email}}
{{- else}}
其他 email:{{$email}}
{{- end}}
{{- end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址:
第一个 email:jerry@qq.com
其他 email:jerry@163.com`
	assert.Equal(t, want, buf.String())
}

func TestTemplateRange(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name:   "Jerry",
		Emails: []string{},
	}
	tmpl, err := template.New("test").Parse(`{{ .Name }}的邮箱地址：
{{- range .Emails}}
{{.}}
{{- else}}
邮箱地址为空
{{- end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址：
邮箱地址为空`
	assert.Equal(t, want, buf.String())
}

func TestTemplateWith(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name: "Jerry",
		Emails: []string{
			"jerry@qq.com",
			"jerry@163.com",
		},
	}
	tmpl, err := template.New("test").Parse(`{{ .Name }}的邮箱地址:
{{- range $idx, $email := .Emails}}
{{$email}}的下标是否等于0: {{eq $idx 0}}
{{- end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址:
jerry@qq.com的下标是否等于0: true
jerry@163.com的下标是否等于0: false`
	assert.Equal(t, want, buf.String())
}

func TestTemplateFuncMap(t *testing.T) {
	type Person struct {
		Name   string
		Emails []string
	}
	p := Person{
		Name: "Jerry",
		Emails: []string{
			"jerry@qq.com",
			"jerry@163.com",
		},
	}

	funcMap := template.FuncMap{
		"inr": func(i int) int { return i + 1 },
	}

	tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{ .Name }}的邮箱地址：
{{range $index, $email := .Emails}}
{{- $index = inr $index}}
{{- printf "%d. %s\n" $index $email }}
{{- end}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, p)
	assert.NoError(t, err)
	want := `Jerry的邮箱地址：
1. jerry@qq.com
2. jerry@163.com
`
	assert.Equal(t, want, buf.String())
}

func TestTemplateDefine(t *testing.T) {
	tmpl, err := template.New("tmp1").Parse(`{{- define "T1"}}我是模板T1{{end}}
{{- define "T2"}}我是模板T2{{end}}
{{- define "T3"}}我是模板T3，嵌套了两个模板: {{template "T1"}} {{template "T2"}}{{end}}
{{- template "T3"}}`)
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, nil)
	assert.NoError(t, err)
	want := `我是模板T3，嵌套了两个模板: 我是模板T1 我是模板T2`
	assert.Equal(t, want, buf.String())
}

func TestTemplateBlock(t *testing.T) {
	const (
		master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
		overlay = `{{define "list"}} {{join . ", "}}{{end}} `
	)
	var (
		funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
	if err != nil {
		log.Fatal(err)
	}
	overlayTmpl, err := template.Must(masterTmpl.Clone()).Parse(overlay)
	if err != nil {
		log.Fatal(err)
	}
	if err := masterTmpl.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
	if err := overlayTmpl.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
}
