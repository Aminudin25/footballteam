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
	"footballteam/match"
	"footballteam/match_result"
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
		&match.Match{},
		&match_result.MatchResult{},
		&match_result.Goal{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate:", err)
	}
	fmt.Println("✅ Database migration completed")

	// Jalankan seeder admin
	seedAdminUser(db)
	seedTeams(db)
	seedPlayers(db)
	seedMatch(db)
	seedMatchResults(db)
	
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

	matchRepository := match.NewRepository(db)
	matchService := match.NewService(matchRepository)
	matchHandler := handler.NewMatchHandler(matchService)

	matchResultRepository := match_result.NewRepository(db)
	matchResultService := match_result.NewService(matchResultRepository, playerService, matchService)
	matchResultHandler := handler.NewMatchResultHandler(matchResultService, playerService)



	router := gin.Default()
	api := router.Group("/api/v1")

	// Public routes
	api.POST("/sessions", userHandler.Login)

	// Teams
	api.GET("/teams", teamHandler.GetTeams)
	api.GET("/teams/:id", teamHandler.GetTeamByID)

	// Players
	api.GET("/players", playerHandler.GetPlayers)
	api.GET("/players/:id", playerHandler.GetPlayerByID)
	api.GET("/players/team/:team_id", playerHandler.GetPlayersByTeam)

	// Matches
	api.GET("/matches", matchHandler.GetMatches)
	api.GET("/matches/:id", matchHandler.GetMatchByID)

	// MatchResults
	api.GET("/match_results", matchResultHandler.GetMatchResults)
	api.GET("/match_results/:id", matchResultHandler.GetMatchResultByID)
	api.GET("/match_results/report", matchResultHandler.GetMatchResultsReport)

	// Protected routes (only admin)
	protected := api.Group("/")
	protected.Use(authMiddleware(authService, userService))

	// Teams
	protected.POST("/teams", teamHandler.CreateTeam)
	protected.PUT("/teams/:id", teamHandler.UpdateTeam)
	protected.DELETE("/teams/:id", teamHandler.DeleteTeam)
	protected.POST("/teams/:id/logo", teamHandler.UploadLogo)

	// Players
	protected.POST("/players", playerHandler.CreatePlayer)
	protected.PUT("/players/:id", playerHandler.UpdatePlayer)
	protected.DELETE("/players/:id", playerHandler.DeletePlayer)

	// Matches
	protected.POST("/matches", matchHandler.CreateMatch)
	protected.PUT("/matches/:id", matchHandler.UpdateMatch)
	protected.DELETE("/matches/:id", matchHandler.DeleteMatch)

	// MatchResults
	protected.POST("/match_results", matchResultHandler.CreateMatchResult)

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

func seedMatch(db *gorm.DB) {
	matches := []match.Match{
		{
			Date:       "2025-10-20",
			Time:       "15:00",
			HomeTeamID: 1,
			AwayTeamID: 2,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Date:       "2025-10-22",
			Time:       "18:30",
			HomeTeamID: 3,
			AwayTeamID: 4,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Date:       "2025-10-25",
			Time:       "20:00",
			HomeTeamID: 2,
			AwayTeamID: 3,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for _, m := range matches {
		var existing match.Match
		err := db.Where("date = ? AND time = ? AND home_team_id = ? AND away_team_id = ?", m.Date, m.Time, m.HomeTeamID, m.AwayTeamID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&m).Error; err != nil {
				log.Printf("❌ Gagal menambahkan match: %+v, error: %v", m, err)
			} else {
				log.Printf("✅ Berhasil menambahkan match: %+v", m)
			}
		}
	}
}

func seedMatchResults(db *gorm.DB) {
	var count int64
	db.Model(&match_result.MatchResult{}).Count(&count)
	if count > 0 {
		fmt.Println("MatchResults already exist, skipping seeder...")
		return
	}

	results := []match_result.MatchResult{
		{
			MatchID:   1,
			HomeScore: 2,
			AwayScore: 1,
			Status:    "Home Menang",
			Goals: []match_result.Goal{
				{
					PlayerID: 1, // ID pemain yang mencetak gol
					TeamID:   1, // Home team
					Minute:   15,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					PlayerID: 2,
					TeamID:   1,
					Minute:   60,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					PlayerID: 4,
					TeamID:   2, // Away team
					Minute:   75,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&results).Error; err != nil {
		log.Println("Failed to seed match results:", err)
		return
	}

	fmt.Println("✅ Seeded match results with goals successfully")
}


