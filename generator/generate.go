package generator

import (
	"strconv"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	contextPkg         = protogen.GoImportPath("context")
	ginPkg             = protogen.GoImportPath("github.com/gin-gonic/gin")
	errPkg             = protogen.GoImportPath("errors")
	metadataPkg        = protogen.GoImportPath("google.golang.org/grpc/metadata")
	deprecationComment = "// Deprecated: Do not use."
)

var methodSets = make(map[string]int)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + "_gin.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by github.com/x-lambda/protoc-gen-gin. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the github.com/x-lambda/protoc-gen-gin package it is being compiled against.")
	g.P("// ", contextPkg.Ident(""))
	g.P("// ", metadataPkg.Ident(""))
	g.P("// ", ginPkg.Ident(""), errPkg.Ident(""))
	g.P()

	for _, service := range file.Services {
		genService(gen, file, g, service)
	}

	return g
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	if s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}

	// HTTP Server
	sd := &service{
		Name:     s.GoName,
		FullName: string(s.Desc.FullName()),
		FilePath: file.Desc.Path(),
	}

	ps := strings.Split(gen.Request.GetParameter(), ",")
	for _, v := range ps {
		fields := strings.Split(v, "=")
		if len(fields) == 2 && strings.Trim(fields[0], " ") == "code" {
			sd.ParamCode = fields[1]
		}
		if len(fields) == 2 && strings.Trim(fields[0], " ") == "err_data" {
			errData, err := strconv.ParseBool(fields[1])
			if err != nil {
				errData = false
			}
			sd.ErrData = errData
		}
	}

	for _, method := range s.Methods {
		sd.Methods = append(sd.Methods, genMethod(method)...)
	}

	text := sd.execute()
	g.P(text)
}

func genMethod(m *protogen.Method) []*method {
	var methods []*method

	// 存在 http rule 配置
	// options
	rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
	if rule != nil && ok {
		for _, bind := range rule.AdditionalBindings {
			methods = append(methods, buildHTTPRule(m, bind))
		}

		methods = append(methods, buildHTTPRule(m, rule))
		return methods
	}

	methods = append(methods, defaultMethod(m))
	return methods
}

func defaultMethod(m *protogen.Method) *method {
	// TODO path
	// $prefix + / + ${package}.${service} + / + ${method}
	// /api/demo.v0.Demo/GetName
	md := buildMethodDesc(m, "POST", "")
	md.Body = "*"
	return md
}

func buildHTTPRule(m *protogen.Method, rule *annotations.HttpRule) *method {
	var path, method string
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = "GET"
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = "PUT"
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = "POST"
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = "DELETE"
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = "PATCH"
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}

	md := buildMethodDesc(m, method, path)
	return md
}

func buildMethodDesc(m *protogen.Method, httpMethod string, path string) *method {
	defer func() { methodSets[m.GoName]++ }()
	reply := "*" + m.Output.GoIdent.GoName
	// if output is stream, reply is array of output
	if m.Desc.IsStreamingServer() {
		reply = "[] " + reply
	}
	md := &method{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: m.Input.GoIdent.GoName,
		Reply:   reply,
		Path:    path,
		Method:  httpMethod,
	}

	md.initPathParams()
	return md
}
