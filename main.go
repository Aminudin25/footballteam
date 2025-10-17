package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"footballteam/user"
)

func main() {
	// Konfigurasi koneksi database
	dsn := "root:@tcp(127.0.0.1:3306)/footballteam?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	fmt.Println("✅ Database connection successful")

	// Auto-migrate: membuat tabel berdasarkan struct
	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatal("❌ Failed to migrate:", err)
	}
	fmt.Println("✅ Database migration completed")

	// Jalankan seeder admin
	seedAdminUser(db)
}

// Seeder admin user default
func seedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&user.User{}).Where("email = ?", "admin@xyz.com").Count(&count)

	if count == 0 {
		// Hash password
		password := "admin123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		admin := user.User{
			Name:         "Admin",
			Email:        "admin@xyz.com",
			PasswordHash: string(hashedPassword),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatal("Failed to seed admin user:", err)
		}

		fmt.Println("✅ Admin user seeded successfully (email: admin@xyz.com, password: admin123)")
	} else {
		fmt.Println("ℹ️ Admin user already exists, skipping seeder.")
	}
}
