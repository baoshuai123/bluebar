package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name        string       `mapstructure:"name"`
	Mode        string       `mapstructure:"mode"`
	Port        int          `mapstructure:"port"`
	Version     string       `mapstructure:"version"`
	StartTime   string       `mapstructure:"start_time"`
	MachineID   int64        `mapstructure:"machine_id"`
	LogConfig   *LogConfig   `mapstructure:"log"`
	MysqlConfig *MysqlConfig `mapstructure:"mysql"`
	RedisConfig *RedisConfig `mapstructure:"redis"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"PoolSize"`
}

func Init() (err error) {
	//配置文件名（不带扩展名）
	//viper.SetConfigName("config")
	//直接指定配置文件位置
	viper.SetConfigFile("./conf/config.yaml")
	//在项目中查找配置文件的路径，可以使用相对路径
	//viper.AddConfigPath(".")
	//设置文件类型，这里是yaml文件 配合远程配置中心使用
	viper.SetConfigType("yaml")
	//查找并读取配置文件
	err = viper.ReadInConfig()
	if err != nil { // 处理读取配置文件的错误
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}
	//把读取的信息反序列化到Conf中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed err: %v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件修改了...", e.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed err: %v\n", err)
		}

	})
	return
}
