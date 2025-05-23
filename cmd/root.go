package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var rootCmd = &cobra.Command{}

// Execute execute root cmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

var cfgFile, projectBase, userLicense string

func init() {
	cobra.OnInitialize(initConfig)

	// 在此可以定义自己的flag或者config设置，Cobra支持持久标签(persistent flag)，它对于整个应用为全局
	// 在StringVarP中需要填写`shorthand`，详细见pflag文档
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/sfs.yaml", "config file (defalut in $HOME/.cobra.yaml)")
	//rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	//rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	//rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")

	// Cobra同样支持局部标签(local flag)，并只在直接调用它时运行
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 使用viper可以绑定flag
	//viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	//viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	//viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	//viper.SetDefault("license", "apache")
	viper.SetDefault("ipfs_ip", "192.168.59.128")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	// 勿忘读取config文件，无论是从cfgFile还是从home文件
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// 找到home文件
		home, err := homedir.Dir()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		// 在home文件夹中搜索以“.cobra”为名称的config
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
	viper.SetConfigType("yaml")
	// 读取符合的环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Can not read config:", viper.ConfigFileUsed())
	}
}
