package setup

import (
	"fmt"

	"github.com/amenabe22/chachata_backend/graph/model"
	_ "github.com/amenabe22/chachata_backend/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupModels() *gorm.DB {

	// https://github.com/jackc/pgx
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password=postgres dbname=chachata_db port=5432 sslmode=disable TimeZone=Asia/Shanghai", // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                                                                                                                 // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})

	// db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=postgres")
	if err != nil {
		fmt.Println(err)
	}
	// defer db.Close()
	// run migrations for this user by default
	// usrs := []model.User{
	// 	{
	// 		Email:    "bdere12345@gmail.com",
	// 		Username: "amenabe",
	// 	},
	// }
	// for _, u := range usrs {
	// 	db.Create(&u)
	// }
	// db.Exec("CREATE DATABASE chachata_db")
	// db.LogMode(true)

	// db.Exec("use chachata_db")
	db.AutoMigrate(
		&model.User{},
		&model.Profile{},
		&model.Devices{},
	)
	// db.Migrator().CreateConstraint(&model.Profile{}, "fk_users_devices")
	// db.Migrator().CreateConstraint(&model.Devices{}), "Profiles")
	// db.Model(&model.Devices{}).Add("cust_id", "customers(cust_id)", "CASCADE", "CASCADE")

	// db.AutoMigrate(
	// 	&model.Users{}, &model.Users{}, &model.Message{},
	// 	&model.Chatroom{}, &model.Usr{}, &model.Company{},
	// 	&model.NickName{})
	// db.AutoMigrate(&model.Usr{})
	// // db.AutoMigrate(&model.User{})
	// db.AutoMigrate(&model.Company{})
	// db.AutoMigrate(&model.Admin{})
	// db.Migrator().CreateConstraint(&model.Profile{}, "UserID")
	// db.Migrator().CreateConstraint(&Users{}), "Detail")
	return db
}
