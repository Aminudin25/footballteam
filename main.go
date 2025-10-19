package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"footballteam/auth"
	"footballteam/handler"
	"footballteam/helper"
	"footballteam/player"
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
		&player.Player{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate:", err)
	}
	fmt.Println("✅ Database migration completed")

	// Jalankan seeder admin
	seedAdminUser(db)
	seedTeams(db)
	seedPlayers(db)
	
	authService := auth.NewService()

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository) 
	userHandler := handler.NewUserHandler(userService, authService)

	teamRepository := team.NewRepository(db)
	teamService := team.NewService(teamRepository)
	teamHandler := handler.NewTeamHandler(teamService)

	playerRepository := player.NewRepository(db)
	playerService := player.NewService(playerRepository)
	playerHandler := handler.NewPlayerHandler(playerService)


	router := gin.Default()
	api := router.Group("/api/v1")

	// Public routes
	//Team
	api.POST("/sessions", userHandler.Login)
	api.GET("/teams", teamHandler.GetTeams)
	api.GET("/teams/:id", teamHandler.GetTeamByID)

	// Player (tanpa login)
    api.GET("/players", playerHandler.GetPlayers)
    api.GET("/players/:id", playerHandler.GetPlayerByID)
    api.GET("/players/team/:team_id", playerHandler.GetPlayersByTeam)

	// Protected routes (only admin)
	protected := api.Group("/")
	protected.Use(authMiddleware(authService, userService))

	// Teams
	protected.POST("/teams", teamHandler.CreateTeam)
	protected.PUT("/teams/:id", teamHandler.UpdateTeam)
	protected.DELETE("/teams/:id", teamHandler.DeleteTeam)
	protected.POST("/teams/:id/logo", teamHandler.UploadLogo)

	// Player (dengan login)
    protected.POST("/players", playerHandler.CreatePlayer)
    protected.PUT("/players/:id", playerHandler.UpdatePlayer)
    protected.DELETE("/players/:id", playerHandler.DeletePlayer)

	api.Static("/uploads", "./uploads")

	router.Run()
	
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Bearer tokentokentoken
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))
		
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

		c.Next()

	}
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

func seedPlayers(db *gorm.DB) {
	var count int64
	db.Model(&player.Player{}).Count(&count)
	if count > 0 {
		fmt.Println("Players already exist, skipping seeder...")
		return
	}

	var garudaFC, pahlawanFC team.Team
	db.Where("name = ?", "Garuda FC").First(&garudaFC)
	db.Where("name = ?", "Pahlawan FC").First(&pahlawanFC)

	// Jika tim tidak ditemukan, skip
	if garudaFC.ID == 0 || pahlawanFC.ID == 0 {
		fmt.Println("❌ Teams not found, skipping player seeder.")
		return
	}

	players := []player.Player{
		// Garuda FC
		{
			Name:         "Rizky Hadi",
			Height:       178,
			Weight:       72,
			Position:     "Penyerang",
			Number:       9,
			TeamID:       garudaFC.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Dimas Putra",
			Height:       180,
			Weight:       75,
			Position:     "Gelandang",
			Number:       10,
			TeamID:       garudaFC.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Andri Saputra",
			Height:       185,
			Weight:       80,
			Position:     "Penjaga Gawang",
			Number:       1,
			TeamID:       garudaFC.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},

		// Pahlawan FC
		{
			Name:         "Budi Santoso",
			Height:       177,
			Weight:       70,
			Position:     "Bertahan",
			Number:       5,
			TeamID:       pahlawanFC.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Yoga Prasetyo",
			Height:       174,
			Weight:       68,
			Position:     "Gelandang",
			Number:       8,
			TeamID:       pahlawanFC.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Ahmad Fadli",
			Height:       182,
			Weight:       78,
			Position:     "Penyerang",
			Number:       11,
			TeamID:       pahlawanFC.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	if err := db.Create(&players).Error; err != nil {
		log.Println("❌ Failed to seed players:", err)
		return
	}

	fmt.Println("✅ Seeded 6 default players successfully")
}
