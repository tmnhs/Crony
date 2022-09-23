package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/utils"
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

var (
	//本地Config自动载入顺序
	autoLoadLocalConfigs = []string{
		ExtensionJson,
		ExtensionYaml,
		ExtensionInI,
	}
)



func LoadConfig(env ,serverName,configFileName string) (*models.Config, error) {
	var c models.Config
	var confPath string
	dir:=fmt.Sprintf("%s/%s/%s",serverName,NameSpace,env)
	for _, registerExt := range autoLoadLocalConfigs {
		confPath = path.Join(dir, configFileName+registerExt)
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

	return &c,nil
}
