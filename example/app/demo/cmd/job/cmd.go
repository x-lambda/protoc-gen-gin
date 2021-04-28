package job

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "job",
	Short: "protoc-gen-gin-example job",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
