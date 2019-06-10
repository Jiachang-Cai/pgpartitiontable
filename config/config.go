package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Name string
	Env  string
}

func Init(cfg, env string) error {
	c := Config{
		Name: cfg,
		Env:  env,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("docs") // 如果没有指定配置文件，则解析默认的配置文件
		if c.Env == "pro" {
			viper.SetConfigName("qcloud")
		} else if c.Env == "stg" {
			viper.SetConfigName("test")
		} else {
			viper.SetConfigName("local")
		}
	}
	viper.SetConfigType("toml") // 设置配置文件格式为YAML
	viper.AutomaticEnv()        // 读取匹配的环境变量
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}
