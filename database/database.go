package database

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {

        database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
                TranslateError: true,
        })

        if err != nil {
                panic(err)
        }

        Db = database

        m, err := migrate.New("file://database/migrations", os.Getenv("DATABASE_URL"))
        if err != nil {
                panic(err)
        }

        if err = m.Up(); err != nil && err != migrate.ErrNoChange {
                panic(err)
        }

        logrus.Info("Completed migration")

}