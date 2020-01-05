package anchor

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster"
	"github.com/ZachiNachshon/anchor/cmd/anchor/completion"
	"github.com/ZachiNachshon/anchor/cmd/anchor/docker"
	"github.com/ZachiNachshon/anchor/cmd/anchor/list"
	"github.com/ZachiNachshon/anchor/cmd/anchor/version"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/installer"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/spf13/cobra"
)

type AnchorCmd struct {
	cobraCmd *cobra.Command
	opts     common.CmdRootOptions
}

var validArgs = []string{"cluster", "completion", "docker", "list", "version"}

func NewCommand() *AnchorCmd {
	var rootCmd = &cobra.Command{
		Use:       "anchor",
		Short:     "Utility for local Docker/Kubernetes development environment",
		Long:      `Utility for local Docker/Kubernetes development environment`,
		ValidArgs: validArgs,
	}

	if err := config.CheckPrerequisites(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := installer.NewGoInstaller(common.ShellExec).Check(); err != nil {
		logger.Fatal(err.Error())
	}

	locator.DirLocator = locator.New()
	if err := locator.DirLocator.Scan(); err != nil {
		logger.Fatal(err.Error())
	}

	extractor.CmdExtractor = extractor.New()

	return &AnchorCmd{
		cobraCmd: rootCmd,
	}
}

func (cmd *AnchorCmd) initFlags() error {
	cmd.cobraCmd.PersistentFlags().BoolVarP(
		&common.GlobalOptions.Verbose,
		"verbose",
		"v",
		common.GlobalOptions.Verbose,
		"anchor <command> -v")

	cmd.cobraCmd.PersistentFlags().BoolVar(
		&common.GlobalOptions.DockerRunAutoLog,
		"auto-log",
		common.GlobalOptions.DockerRunAutoLog,
		"anchor docker run <image> --auto-log=false")

	cmd.cobraCmd.PersistentFlags().SortFlags = false
	return nil
}

func (cmd *AnchorCmd) initSubCommands() error {

	//cobra.EnableCommandSorting = false

	// List Commands
	cmd.cobraCmd.AddCommand(list.NewCommand(&cmd.opts).GetCobraCmd())

	// Docker Commands
	cmd.cobraCmd.AddCommand(docker.NewCommand(&cmd.opts).GetCobraCmd())

	// Kubernetes Commands
	cmd.cobraCmd.AddCommand(cluster.NewCommand(&cmd.opts).GetCobraCmd())

	// Admin
	cmd.cobraCmd.AddCommand(version.NewCommand(&cmd.opts).GetCobraCmd())

	// Auto completion
	cmd.cobraCmd.AddCommand(completion.NewCommand(cmd.cobraCmd, &cmd.opts).GetCobraCmd())

	return nil
}

func (cmd *AnchorCmd) Execute() {

	if err := cmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := cmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := cmd.cobraCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}

func (cmd *AnchorCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func Main() {
	NewCommand().Execute()
}

//func init() {
//	cobra.OnInitialize(initConfig)
//	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
//	rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
//	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
//	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
//	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
//	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
//	viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
//	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
//	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
//	viper.SetDefault("license", "apache")
//}
//
//func initConfig() {
//	// Don't forget to read config either from cfgFile or from home directory!
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		// Search config in home directory with name ".cobra" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".cobra")
//	}
//
//	if err := viper.ReadInConfig(); err != nil {
//		fmt.Println("Can't read config:", err)
//		os.Exit(1)
//	}
//}
