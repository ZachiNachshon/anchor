package cmd

import (
	"fmt"
	"github.com/kit/cmd/docker"
	"github.com/kit/cmd/logger"
	"github.com/kit/cmd/types"
	"github.com/spf13/cobra"
)

var Verbose bool
var Source string

type CmdRoot struct {
	cobraCmd *cobra.Command
	opts     common.CmdRootOptions
}

func NewCmdRoot() *CmdRoot {
	var rootCmd = &cobra.Command{
		Use:   "kit",
		Short: "Short message",
		Long:  `Long multi line message`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("I'm here at root")
		},
	}

	return &CmdRoot{
		cobraCmd: rootCmd,
	}
}

func (root *CmdRoot) initFlags() error {
	root.cobraCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}

func (root *CmdRoot) initSubCommands() error {

	// Docker Commands
	root.cobraCmd.AddCommand(docker.NewDockerCmd(&root.opts).GetCobraCmd())

	// Kubernetes Commands
	//root.cobraCmd.AddCommand(kubernetes.NewKindCmd(&root.opts).GetCobraCmd())

	// Admin
	root.cobraCmd.AddCommand(NewVersionCmd(&root.opts).GetCobraCmd())

	return nil
}

func (root *CmdRoot) Execute() {

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
//	rootCmd.AddCommand(VersionCmd)
//}

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
