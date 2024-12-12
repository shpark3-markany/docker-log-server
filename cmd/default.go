package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use: "default",
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()

    

		e.Logger.Fatal(e.Start(":1220"))
	},
}

func init() {
	rootCmd.AddCommand(defaultCmd)
}
