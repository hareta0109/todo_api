package main

import (
	"fmt"
	"log"
	"time"
	"todo_api/internal/adapter/outbound/mysql/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func p[T any](value T) *T {
	ptr := new(T)
	*ptr = value
	return ptr
}

func seeds(db *gorm.DB) error {
	companies := []model.Company{
		{
			ID:   1,
			Name: "管理会社",
		},
		{
			ID:   2,
			Name: "利用会社",
		},
	}

	users := []model.Auth{
		{
			ID:   1,
			Name: "管理会社の管理者",
			// super-admin
			Hash:      "11c520bc9f1f460f7581c2514c24e652f96efbe796e8287613825bbed2ae6618",
			Role:      "EDITOR",
			UserType:  "ADMIN",
			CompanyID: 1,
		},
		{
			ID:   2,
			Name: "管理会社のユーザ",
			// super-user
			Hash:      "d761b775020614c19fc5f717c4d760eec7e045fdce450fd26d4e59a36024db3b",
			Role:      "EDITOR",
			UserType:  "NORMAL",
			CompanyID: 1,
		},
		{
			ID:   3,
			Name: "利用会社の管理者",
			// normal-admin
			Hash:      "43829479f47c75412898555caf8896b67830e7af9e81b9a56d29dd9e8bec507d",
			Role:      "EDITOR",
			UserType:  "ADMIN",
			CompanyID: 2,
		},
		{
			ID:   4,
			Name: "利用会社の一般編集者",
			// normal-user
			Hash:      "45cefa501529bb09955fb45f28207e99ab0e6adb85d0939e2c95572e632a8954",
			Role:      "EDITOR",
			UserType:  "NORMAL",
			CompanyID: 2,
		},
		{
			ID:   5,
			Name: "利用会社の一般閲覧者",
			// normal-user
			Hash:      "45cefa501529bb09955fb45f28207e99ab0e6adb85d0939e2c95572e632a8954",
			Role:      "VIEWER",
			UserType:  "NORMAL",
			CompanyID: 2,
		},
	}

	tasks := []model.Task{
		{
			ID:               1,
			Title:            "最初のタスク",
			Detail:           p("タスク詳細"),
			Status:           "NEW",
			Visibility:       "COMPANY",
			PersonInChargeID: p(uint64(4)),
			LimitDate:        p(time.Now().AddDate(0, 1, 0)),
			CreatorID:        3,
			UpdatorID:        3,
			CreateAt:         time.Now(),
			UpdateAt:         time.Now(),
		},
	}

	if err := db.Create(&companies).Error; err != nil {
		fmt.Printf("%+v", err)
	}

	if err := db.Create(&users).Error; err != nil {
		fmt.Printf("%+v", err)
	}

	if err := db.Create(&tasks).Error; err != nil {
		fmt.Printf("%+v", err)
	}

	return nil
}

func openConnection() *gorm.DB {
	dsn := "user:pass@tcp(127.0.0.1:3306)/task_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Couldn't establish database connection: %s", err)
	}
	return db
}

func main() {
	db := openConnection()
	if err := seeds(db); err != nil {
		fmt.Printf("%+v", err)
		return
	}
}
