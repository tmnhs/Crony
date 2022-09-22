package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/utils"
	"os"
	"path"
)

const (
	//ExtensionJson json配置后缀
	ExtensionJson = ".json"
	//ExtensionYaml yaml配置后缀
	ExtensionYaml = ".yaml"
	//ExtensionInI ini配置后缀
	ExtensionInI = ".ini"

	NameSpace = "conf"
)

var c models.Config
var (
	//本地Config自动载入顺序
	autoLoadLocalConfigs = []string{
		ExtensionJson,
		ExtensionYaml,
		ExtensionInI,
	}
)

var ConfOptions struct {
	flags.Options
	Environment string `short:"e" long:"env" description:"Use crony-server environment"`
	Version     bool   `short:"v" long:"version"  description:"Show crony-server version"`
	//todo
	EnablePProfile    bool   `short:"p" long:"enable-pprof"  description:"enable pprof"`
	EnableHealthCheck bool   `short:"a" long:"enable-health-check"  description:"enable health check"`
	ConfigFileName    string `short:"c" long:"config" description:"Use coa-server config file name" default:"main"`
	EnableDevMode     bool   `short:"m" long:"enable-dev-mode"  description:"enable dev mode"`
}

func LoadConfig(profile string) error {
	var parser = flags.NewParser(&ConfOptions, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		return err
	}
	fmt.Println("confOptions is :", ConfOptions)
	var confPath string
	configFileName := ConfOptions.ConfigFileName
	if configFileName == "" {
		configFileName = "main"
	}
	for _, registerExt := range autoLoadLocalConfigs {
		confPath = path.Join(profile+"/"+NameSpace, configFileName+registerExt)
		if utils.Exists(confPath) {
			break
			//return NewConfig(env, namespace, configFileName+registerExt)
		}
	}
	fmt.Println("confPath is :", confPath)
	v := viper.New()
	v.SetConfigFile(confPath)
	ext := utils.Ext(confPath)
	v.SetConfigType(ext)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&c); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&c); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("config is :%#v", c)
	return nil
}
