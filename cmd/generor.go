package cmd

import (
	"fmt"

	"github.com/jaxenlau/pagenor-go/services"
	"github.com/spf13/cobra"
)

// serverCmd represents the http api server command
var serverCmd = &cobra.Command{
	Use:   "gen",
	Short: "create page",
	Run: func(cmd *cobra.Command, args []string) {
		opts := loadApplicationOptions()

		fmt.Println(opts.Pagenor)

		pg := services.NewPagenor(&opts.Pagenor)

		err := pg.Generate()
		handleInitError("generate", err)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
