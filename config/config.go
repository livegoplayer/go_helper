package config

import (
	"fmt"
	"github.com/livegoplayer/go_helper/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const PRODUCTION_ENV = "PRO"
const DEVELOPMENT_ENV = "DEV"
const ENV_PREFIX = "US"

func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return utils.PathToCommon(dir)
}

//这是一个设置环境变量的方法
func SetEnv(envName string, envVal string) {
	viper.SetEnvPrefix(ENV_PREFIX)
	err := viper.BindEnv(envName)
	if err != nil {
		panic("环境变量绑定出错")
	}

	//之后可以通过 viper.get 获取，无需prefix
	err = os.Setenv(ENV_PREFIX+"_"+envName, envVal)
	if err != nil {
		panic("设置环境变量出错")
	}
}

func LoadEnv() {
	//首先加载文件目录下的.env文件,只有文件加载失败才会出错，文件不存在不会
	if err := godotenv.Load(utils.PathToCommon(GetCurrentPath() + "/.env")); err != nil {
		panic("文件加载出错:" + err.Error())
	}

	//默认是dev环境
	viper.SetDefault("ENV", DEVELOPMENT_ENV)
	//获取当前系统路径,读取.env文件的值
	viper.SetDefault("ROOT", utils.PathToCommon(GetCurrentPath()))

	//加载完成之后可以使用os.Getenv()方法来获取对应的env，这里需要把他们加载到viper中，即在viper中注册一下 todo 添加更多
	SetEnv("APP_DEBUG", os.Getenv("APP_DEBUG"))
	SetEnv("ENV", os.Getenv("ENV"))
	SetEnv("DEV_CONFIG_PATH", os.Getenv("DEV_CONFIG_PATH"))
	SetEnv("PRO_CONFIG_PATH", os.Getenv("PRO_CONFIG_PATH"))

	//设置配置文件目录
	var configFilePath string
	switch viper.GetString("ENV") {
	case PRODUCTION_ENV:
		configFilePath = viper.GetString("ROOT") + viper.GetString("PRO_CONFIG_PATH")
	case DEVELOPMENT_ENV:
		configFilePath = viper.GetString("ROOT") + viper.GetString("DEV_CONFIG_PATH")
	default:
		configFilePath = viper.GetString("ROOT") + viper.GetString("PRO_CONFIG_PATH")
	}

	viper.Set("CONFIG_FILE_PATH", utils.PathToCommon(configFilePath))
	viper.AddConfigPath(viper.GetString("CONFIG_FILE_PATH"))
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
