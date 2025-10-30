package model

type Client struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	DBName     string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	IsActive   bool
}
