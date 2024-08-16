package database

import (
	"fmt"
	"splitz/config"
	"splitz/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cnf *config.EnvModel) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cnf.DbHost, cnf.DbUser, cnf.DbName, cnf.DbPort, cnf.DbPassword)
	db, err := gorm.Open(postgres.Open(psqlInfo))

	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Group{})
	db.AutoMigrate(&domain.UserGroup{})
	db.AutoMigrate(&domain.Request{})
	db.AutoMigrate(&domain.UserRequest{})

	return db, err

}
