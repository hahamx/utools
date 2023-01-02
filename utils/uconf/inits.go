package uconf

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hahamx/utools/models/config"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

var (
	ConfigEnv  = "GCONFIG"
	ConfigFile = "config.yaml"
	LocalFile  = "./config.yaml"
)

type CfIniter struct {
	*viper.Viper
	MysqlDsn string
	All      *config.Server
}

// InitConfig优先级: 命令行 > 环境变量 > 默认值
func InitConfig(path ...string) (*viper.Viper, config.Server) {
	var local string
	var newGconfig = config.Server{}

	if len(path) == 0 {
		flag.StringVar(&local, "c", "", "choose config file.")
		flag.Parse()
		fmt.Printf("flags:%v\n", flag.Args())
		if local == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(ConfigFile); configEnv == "" { // 判断 internal.ConfigEnv 常量存储的环境变量是否为空
				local = LocalFile
				fmt.Printf("服务模式:%s,配置:%s\n", gin.EnvGinMode, LocalFile)

			} else { // internal.ConfigEnv 常量存储的环境变量不为空 将值赋值于config
				local = configEnv
				fmt.Printf("服务模式:%s,配置:%s\n", ConfigEnv, local)
			}
		} else { // 命令行参数不为空 将值赋值于config
			fmt.Printf("命令行的-c参数值,配置为:%s\n", local)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		local = path[0]
		fmt.Printf("服务模式:%s ,配置:%s\n", gin.EnvGinMode, local)
	}

	v := viper.New()
	v.SetConfigFile(local)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %v \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&newGconfig); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&newGconfig); err != nil {
		fmt.Println(err)
	}
	newGconfig.AutoCode.Root, _ = filepath.Abs("..")

	return v, newGconfig
}
