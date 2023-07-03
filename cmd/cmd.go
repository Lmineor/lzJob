package cmd

import (
	goflag "flag"
	"fmt"
	"github.com/Lmineor/lzJob/api"
	"github.com/Lmineor/lzJob/config"
	"github.com/Lmineor/lzJob/context"
	"github.com/Lmineor/lzJob/pqueue"
	"github.com/Lmineor/lzJob/store"
	"github.com/Lmineor/lzJob/task"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			visitAllFlags(cmd.Flags())
			fmt.Println("v1.0")
		},
	}
	return cmd
}

func RootCmd() *cobra.Command {
	flags := NewFlags()
	cmd := &cobra.Command{
		Use:   "lj",
		Short: "lj aka lzJob designed to manage timed jobs",
		Long: `
lj aka lzJob designed to manage timed jobs.
`,
	}
	flags.AddFlags(cmd.PersistentFlags())
	cmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
	cmd.AddCommand(serverCmd(flags))
	return cmd
}

func serverCmd(f *Flags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "start node server",
		Long: `
Run server.
`,
		Run: func(cmd *cobra.Command, args []string) {
			visitAllFlags(cmd.Flags())
			run(f)
		},
	}
	return cmd
}
func initCfg(f *Flags) *config.Config {
	var cfg config.Config
	if f.ConfigFile == "" {
		dbCfg := config.Mysql{
			Path:         f.DbUrl,
			Port:         f.DbPort,
			Username:     f.DbUserName,
			Password:     f.DbPassword,
			Db:           f.Db,
			Config:       f.DbConfig,
			MaxIdleConns: f.MaxIdleConns,
			MaxOpenConns: f.MaxOpenConns,
		}
		sysCfg := config.Server{
			ServerAddr: f.Addr,
			ListenPort: f.Port,
		}
		cfg.Server = sysCfg
		cfg.Mysql = dbCfg
	} else {
		cfg = *config.InitConfig(f.ConfigFile)
	}
	return &cfg
}

func run(f *Flags) {
	cfg := initCfg(f)
	db := store.InitMysql(&cfg.Mysql)
	pq := pqueue.NewPriorityQueue()

	ctx := context.New(cfg, db, pq)

	if ctx.DB != nil {
		registerTables(ctx)
		db, _ := ctx.DB.DB()
		// 程序结束前关闭数据库链接
		defer db.Close()
	}
	go task.Trigger(ctx)
	srv := api.NewServer(ctx)
	srv.ListenAndServe()

}
func visitAllFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		klog.Infof("Flag: --%s=%q", flag.Name, flag.Value)
	})
}

func registerTables(ctx context.LZContext) {
	err := ctx.DB.AutoMigrate(
		&store.Task{},
		&store.TaskResult{},
	)
	if err != nil {
		klog.Fatalf("register table failed %s", err)
	}
	klog.Info("register table success")
}
