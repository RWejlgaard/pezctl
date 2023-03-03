package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var Verbose bool
var Logger *zap.SugaredLogger

func init() {
    rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "logs a lot")

    logConfig, _ := zap.NewProduction()
    Logger = logConfig.Sugar()
    if Verbose {
        Logger.Infow("Logging enabled")
    }
}

var rootCmd = &cobra.Command{
    Use: "pezctl",
    Short: "PezCtl is for managing pez-sh",

    Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
