package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/hcolde/reviewer-helper/conf"
)

type Mysql struct {
	DB *gorm.DB
}

func (m *Mysql) New() error {
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Conf.Mysql.User,
		conf.Conf.Mysql.Password,
		conf.Conf.Mysql.Type,
		conf.Conf.Mysql.Path,
		conf.Conf.Mysql.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	m.DB = db
	return nil
}
