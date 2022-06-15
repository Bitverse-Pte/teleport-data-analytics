package cmd

import (
	"github.com/teleport-network/teleport-data-analytics/jobs"
	"github.com/teleport-network/teleport-data-analytics/version"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"

	"github.com/teleport-network/teleport-data-analytics/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "teleport-data-analytics",
		Short: "",
		Run:   func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
	}
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start teleport bridge backend.",
		Run:   func(cmd *cobra.Command, args []string) { Run() },
	}
	versionCmd = version.NewVersionCommand()
)

func init() {
	startCmd.Flags().StringVarP(&config.LocalConfig, "config", "c", "", "")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func Run() {
	cfg := config.LoadConfigs()
	scheduler := gocron.NewScheduler(time.UTC)
	pktSrv := jobs.NewPacketService(scheduler, cfg)
	//go pktSrv.PktPool.UpdateNoMutilIdData()
	pktSrv.PktPool.SyncToDB(scheduler, cfg.SyncEnable)
	scheduler.StartBlocking()
}
