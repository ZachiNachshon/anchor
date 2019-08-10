package cmd

import (
	"github.com/ZachiNachshon/anchor/cmd/docker"
	"github.com/ZachiNachshon/anchor/cmd/kubernetes"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	cobraCmd *cobra.Command
	opts     common.CmdRootOptions
}

var validArgs = []string{"docker", "kubernetes", "list", "version", "completion"}

func NewCmdRoot() *RootCmd {
	var rootCmd = &cobra.Command{
		Use:   "anchor",
		Short: "Utility for local Docker/Kubernetes development environment",
		Long:  `Utility for local Docker/Kubernetes development environment`,
		//BashCompletionFunction: config.AutoCompletionFuncZsh,
		//ValidArgs:              validArgs,
	}

	if err := config.CheckPrerequisites(); err != nil {
		logger.Fatal(err.Error())
	}

	locator.DirLocator = locator.New()
	if err := locator.DirLocator.Scan(); err != nil {
		logger.Fatal(err.Error())
	}

	extractor.CmdExtractor = extractor.New()

	return &RootCmd{
		cobraCmd: rootCmd,
	}
}

func (root *RootCmd) initFlags() error {
	root.cobraCmd.PersistentFlags().BoolVarP(
		&common.GlobalOptions.Verbose,
		"verbose",
		"v",
		common.GlobalOptions.Verbose,
		"anchor <command> -v")

	root.cobraCmd.PersistentFlags().SortFlags = false
	return nil
}

func (root *RootCmd) initSubCommands() error {

	//cobra.EnableCommandSorting = false

	// List Commands
	root.cobraCmd.AddCommand(NewListCmd(&root.opts).GetCobraCmd())

	// Docker Commands
	root.cobraCmd.AddCommand(docker.NewDockerCmd(&root.opts).GetCobraCmd())

	// Kubernetes Commands
	root.cobraCmd.AddCommand(kubernetes.NewKubernetesCmd(&root.opts).GetCobraCmd())

	// Admin
	root.cobraCmd.AddCommand(NewVersionCmd(&root.opts).GetCobraCmd())

	// Auto completion
	root.cobraCmd.AddCommand(NewCompletionCmd(root, &root.opts).GetCobraCmd())

	return nil
}

func (root *RootCmd) Execute() {

	if err := root.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := root.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := root.cobraCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
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
