package postgres

import (
	"auth/src/config"
	"auth/src/refreshSession"
	"auth/src/user"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < config.DbRetries; i++ {
		db, err = gorm.Open(postgres.Open(config.PgConnStr), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, err
	}
	user.Init(db)
	refreshSession.Init(db)

	return db, nil
}
