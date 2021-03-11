package cmd

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()
		e.Use(middleware.Logger())
		e.Any("/test", func(c echo.Context) error {
			for name, values := range c.Request().Header {
				if !strings.HasPrefix(strings.ToLower(name), "httprepeater-") {
					continue
				}
				for _, value := range values {
					c.Response().Header().Add(name, value)
				}
			}

			defer c.Request().Body.Close()
			content, _ := ioutil.ReadAll(c.Request().Body)

			return c.JSON(http.StatusOK, map[string]interface{}{
				"method":  c.Request().Method,
				"uri":     c.Request().URL.RequestURI(),
				"body":    string(content),
				"headers": c.Request().Header,
			})
		})
		port := os.Getenv("PORT")
		addr := ":" + port
		if port == "" {
			addr = "127.0.0.1:1323"
		}
		e.Logger.Fatal(e.Start(addr))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
