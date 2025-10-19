package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"footballteam/auth"
	"footballteam/handler"
	"footballteam/team"
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
	err = db.AutoMigrate(
		&user.User{},
		&team.Team{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate:", err)
	}
	fmt.Println("✅ Database migration completed")

	// Jalankan seeder admin
	seedAdminUser(db)
	seedTeams(db)
	
	authService := auth.NewService()

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository) 
	userHandler := handler.NewUserHandler(userService, authService)

	teamRepository := team.NewRepository(db)
	teamService := team.NewService(teamRepository)
	teamHandler := handler.NewTeamHandler(teamService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/sessions", userHandler.Login)

	api.GET("/teams", teamHandler.GetTeams)
	api.GET("/teams/:id", teamHandler.GetTeamByID)
	api.POST("/teams", teamHandler.CreateTeam)
	api.PUT("/teams/:id", teamHandler.UpdateTeam)
	api.DELETE("/teams/:id", teamHandler.DeleteTeam)
	api.POST("teams/logo/:id", teamHandler.UploadLogo)

	api.Static("/uploads", "./uploads")

	router.Run()
	

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

func seedTeams(db *gorm.DB) {
	var count int64
	db.Model(&team.Team{}).Count(&count)
	if count > 0 {
		fmt.Println("Teams already exist, skipping seeder...")
		return
	}

	teams := []team.Team{
		{
			Name:        "Garuda FC",
			Logo:        "",
			YearFounded: 1998,
			Address:     "Jl. Merdeka No. 1",
			City:        "Jakarta",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Pahlawan FC",
			Logo:        "",
			YearFounded:  2005,
			Address:     "Jl. Pahlawan No. 7",
			City:        "Surabaya",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.Create(&teams).Error; err != nil {
		log.Println("Failed to seed teams:", err)
		return
	}

	fmt.Println("✅ Seeded 2 default teams successfully")
}
