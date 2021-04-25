package main

import (
	"github.com/spf13/cobra"
	"github.com/x-lambda/protoc-gen-gin-example/cmd/help"
	"github.com/x-lambda/protoc-gen-gin-example/cmd/job"
	"github.com/x-lambda/protoc-gen-gin-example/cmd/server"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "",
		Short: "protoc-gen-gin-example",
		Long:  "protoc-gen-gin example",
		Run: func(cmd *cobra.Command, args []string) {
			help.Usage()
		},
	}

	rootCmd.AddCommand(
		help.Cmd,
		server.Cmd,
		job.Cmd,
	)

	rootCmd.Execute()
}
