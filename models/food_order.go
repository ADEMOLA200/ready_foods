package models

import "time"

type MenuItem struct {
    ID      uint   `gorm:"primaryKey" json:"id"`
    Name    string `json:"name"`
    Price   int    `json:"price"`
    OrderID uint   `json:"-"`
}


type FoodOrder struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    UserPhone   string     `json:"user_phone"`
    Name        string     `json:"name"`
    Address     string     `json:"address"`
    MenuItems   []MenuItem `gorm:"foreignKey:OrderID" json:"menu_items"`
    TotalItems  int        `json:"total_items"`
    TotalPay    float64    `json:"total_pay"`
    CreateDtm   time.Time  `json:"create_dtm"`
}
