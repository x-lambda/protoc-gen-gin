package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed template.go.tpl
var tpl string

// rpc GetDemoName(*Req, *Resp)
type method struct {
	Name    string // GetDemoName
	Num     int    // 一个 rpc 方法可以对应多个http请求
	Request string // *Req
	Reply   string // *Resp

	// http rule
	Path         string
	Method       string
	Body         string
	ResponseBody string
}

func (m *method) HandlerName() string {
	return fmt.Sprintf("%s_%d", m.Name, m.Num)
}

// HasPathParams 是否包含路由参数
func (m *method) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}

	return false
}

// initPathParams 转换参数路由 {xx} --> :xx
func (m *method) initPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, p := range paths {
		if p != "" && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			paths[i] = ":" + p[1:len(p)-1]
		}
	}

	m.Path = strings.Join(paths, "/")
}

type service struct {
	Name     string
	FullName string
	FilePath string

	Methods   []*method
	MethodSet map[string]*method
}

func (s *service) execute() string {
	if s.MethodSet == nil {
		s.MethodSet = make(map[string]*method, len(s.Methods))

		for _, m := range s.Methods {
			m := m // TODO ?
			s.MethodSet[m.Name] = m
		}
	}

	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}

	return buf.String()
}

func (s *service) InterfaceName() string {
	return s.Name + "HTTPServer"
}
