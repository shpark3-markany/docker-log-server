package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markany-inc/docker-log-server/controllers"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {

		controllers.MakeClient()

		e := echo.New()
		e.Use(middleware.Logger())
		e.GET("/health", healthCheck)
		e.GET("/version", versionHandler)

		logGroup := e.Group("/api/logs")
		logGroup.GET("/get", controllers.GetLog)
		logGroup.GET("/stream", controllers.GetLog)
		
		e.Logger.Fatal(e.Start(":1220"))
	},
}

func versionHandler(c echo.Context) error {
	return c.String(200, "v1.0.0")
}

func healthCheck(c echo.Context) error {
	return c.String(200, "OK")
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
