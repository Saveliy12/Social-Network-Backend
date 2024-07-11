package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Saveliy12/prod2/internal/api"
	"github.com/Saveliy12/prod2/internal/config"
	"github.com/Saveliy12/prod2/internal/database"
	logger "github.com/Saveliy12/prod2/internal/logger"
	"github.com/Saveliy12/prod2/internal/service"
	"github.com/Saveliy12/prod2/internal/service/tokenmanager"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func main() {

	logger.InitLogger()
	log := logger.GetLogger()

	cfg, err := config.New()
	if err != nil {
		log.Logger.Fatal(err)
	}

	log.Debugf("CONFIG: %+v\n", cfg)

	// Инициализация базы данных
	db := initDB(cfg)

	// database.CreateTables(db) // Создание необходимых таблиц в базе данных

	// Инициализация репозиториев
	userRepository := database.NewUserRepository(db)

	// Инициализация менеджера работы с токенами
	tokenManager, _ := tokenmanager.NewManager("your-signing-key") // подгружать из окружения хз

	// Инициализация сервисов
	// вынести в константы ttl
	authService := service.NewAuthService(tokenManager, userRepository, time.Hour*1, time.Hour*24*30)

	authHandler := api.NewAuthHandler(authService)

	// Инициализация роутеров
	r := gin.Default()

	// Эндпоинты для аутентификации и регистрации
	r.POST("/register", authHandler.RegisterUserHandler)
	r.POST("/login", authHandler.LoginUserHandler)
	r.POST("/refresh", authHandler.RefreshTokenHandler)
	// r.POST("/logout", authHandler.LogoutHandler)

	// Защищенные маршруты
	protected := r.Group("/protected")
	protected.Use(api.AuthMiddleware(tokenManager))
	protected.GET("/profile", authHandler.ProtectedProfileHandler)

	// Запускаем сервер на порту :8080
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Logger.Fatal("Error starting server: ", err)
	}

}

// initRoutes инициализирует все маршруты
func initRouter(authHandler api.AuthHandlerInterface) *gin.Engine {
	// Создаем новый экземпляр Gin
	r := gin.Default()

	// Инициализируем маршруты
	r.GET("/api/ping", api.PingHandler)
	r.POST("/api/register", authHandler.RegisterUserHandler)
	r.POST("/api/login", authHandler.LoginUserHandler)
	r.GET("/verify", authHandler.VerifyTokenHandler)

	return r
}

func initDB(cfg *config.Config) *sqlx.DB {

	// Подключение к базе данных
	info := fmt.Sprintf("user=%v dbname=%v sslmode=%v password=%v", cfg.DB.User, cfg.DB.Name, cfg.DB.SSLMode, cfg.DB.Password)
	db, err := sqlx.Connect("postgres", info)
	if err != nil {
		log.Logger.Fatal("Error connecting to the database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Logger.Fatal("Error pinging the database: ", err)
	}

	return db
}
