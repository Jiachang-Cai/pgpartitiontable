package main

import (
	"os"
	"runtime"
	"flag"

	"pgpartitiontable/config"
	"pgpartitiontable/models"
	"pgpartitiontable/cache"
	"pgpartitiontable/scripts"

	"github.com/robfig/cron"


)

var (
	tomlFile string
)

func init() {
	flag.StringVar(&tomlFile, "config", "docs/local.toml", "配置文件")
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	flag.Parse()
	// init config
	if err := config.Init(tomlFile, os.Getenv("MODE")); err != nil {
		panic(err)
	}
	// init db
	models.DB.Init()
	defer models.DB.Close()
	// init cache
	cache.RedisConn.Init()
	c := cron.New()
	// 每天 0点
	c.AddFunc("0 0 0 * * *", func() {
		scripts.TimeStamp(3)
		scripts.DateTime(3)
	})
	c.Start()
	select {}

}
