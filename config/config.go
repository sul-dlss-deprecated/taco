package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(cfgFile string) {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".taco") // name of config file (without extension)
		viper.AddConfigPath("$HOME") // adding home directory as first search path
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Printf("Error no config file: %s", err)
	}
}

func relativePath(basedir string, path *string) {
	p := *path
	if p != "" && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}
