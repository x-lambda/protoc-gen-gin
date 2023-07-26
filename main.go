package main

import (
	"flag"

	"github.com/x-lambda/protoc-gen-gin/generator"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var flags flag.FlagSet

	var paramErrCode string
	var errData bool
	flags.StringVar(&paramErrCode, "code", "400", "the code of response data when param parse error")
	flags.BoolVar(&errData, "err_data", false, "whether to return data when error occurs")

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			generator.GenerateFile(gen, f)
		}

		return nil
	})
}
