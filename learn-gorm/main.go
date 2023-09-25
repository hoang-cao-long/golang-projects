package main

//TODO
//FIXME
//HACK
//XXX
//BUG
import (
	"log"

	"encoding/json"
	"fmt"

	"os"
	"time"

	"github.com/hoang-cao-long/learn-gorm/config"
	"github.com/hoang-cao-long/learn-gorm/database"
	"github.com/hoang-cao-long/learn-gorm/model"
	"github.com/spf13/viper"
)

func InitConfig() (config.Config, error) {
	var configFile config.Config

	home, _ := os.Getwd()

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config.Config{}, err
	} else {
		viper.Unmarshal(&configFile)
		return configFile, nil
	}
}

func main() {
	config, err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.InitGORM(config)
	if err != nil {
		log.Println("Cannot connect to database")
	}

	now := time.Now().UTC()

	item := model.TodoItem{
		Id:          1,
		Title:       "This is item 1",
		Description: "This is item 1",
		Status:      "Doing",
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	jsonData, err := json.Marshal(item)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))
	jsonStr := "{\"id\":1,\"title\":\"This is item 1\",\"description\":\"This is item 1\",\"status\":\"Doing\",\"created_at\":\"2023-08-20T15:11:40.892347638Z\",\"updated_at\":null}"
	var item2 model.TodoItem

	if err := json.Unmarshal([]byte(jsonStr), &item2); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item2)

	db.AutoMigrate(&model.Category{})

	db.Create(&model.Category{Name: "long"})

	// var cate []model.Category
	// db.First(&cate, 1)
	// db.Find(&cate)

	// fmt.Println(cate)

	// // Update - update product's price to 200
	// db.Model(&product).Update("Price", 200)
	// // Update - update multiple fields
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// // Delete - delete product
	// db.Delete(&product, 1)
}
