package models

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/go-xweb/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	Self *gorm.DB
	Read *gorm.DB
}

type databaseConfig struct {
	Host     string
	Port     int
	User     string
	PassWord string
	DBName   string
}

var DB *Database

func openDB(c *databaseConfig) *gorm.DB {

	config := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.PassWord, c.DBName)
	db, err := gorm.Open("postgres", config)
	if err != nil {
		log.Errorf("Database connection failed. Database name: %s err:%v", c.DBName, err)
	}
	// set for db connection
	setupDB(db)
	return db
}

func setupDB(db *gorm.DB) {
	//db.LogMode(viper.GetBool("gormlog"))
	db.DB().SetMaxOpenConns(30) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0)  // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

func InitSelfDB() *gorm.DB {

	var c databaseConfig
	if err := viper.UnmarshalKey("dbservers.db", &c); err != nil {
		log.Errorf("unable to decode database into struct. err:%v", err)
	}
	return openDB(&c)
}



func (db *Database) Init() {
	DB = &Database{
		Self: InitSelfDB(),
	}
}

func (db *Database) Close() {
	DB.Self.Close()
	DB.Read.Close()
}
