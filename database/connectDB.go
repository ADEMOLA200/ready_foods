package database

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/ADEMOLA200/danas-food/models"
)

var DB *gorm.DB

func ConnectDB() error {
    dsn := "root:rootroot@tcp(127.0.0.1:3306)/ready_food?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }
    
    DB = db

    // AutoMigrate models
    if err := DB.AutoMigrate(&models.FoodOrder{}, &models.MenuItem{}, &models.User{}); err != nil {
        return err
    }

    return nil
}
