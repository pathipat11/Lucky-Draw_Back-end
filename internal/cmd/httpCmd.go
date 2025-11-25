package cmd

import (
	"app/app/routes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
)

func HttpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run server on HTTP protocol",
		Run: func(cmd *cobra.Command, args []string) {
			r := gin.Default()
			routes.Router(r)

			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			r.Run("0.0.0.0:" + port)
		},
	}
}
