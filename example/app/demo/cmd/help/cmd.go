package help

import "github.com/spf13/cobra"

func Usage() {}

var Cmd = &cobra.Command{
	Use:   "help",
	Short: "protoc-gen-gin-example help",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		Usage()
	},
}
