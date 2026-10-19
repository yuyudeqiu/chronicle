package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/config"
	"github.com/yuyudeqiu/chronicle/internal/handler"
	"github.com/yuyudeqiu/chronicle/internal/model"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the chronicle HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Data directory: %s", config.Load())

		r := newAPIRouter()

		// Start server
		log.Println("Starting server on :8080...")
		if err := r.Run(":8080"); err != nil {
			return fmt.Errorf("start server: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func newAPIRouter() *gin.Engine {
	r := gin.Default()
	handler.RegisterRoutes(r)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.ErrorResp(http.StatusNotFound, "route not found"))
	})
	return r
}
