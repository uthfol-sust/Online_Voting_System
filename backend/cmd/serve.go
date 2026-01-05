package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pollvoting/pkg/config"
	"pollvoting/pkg/controllers"
	"pollvoting/pkg/database"
	"pollvoting/pkg/repositories"
	"pollvoting/pkg/routes"
	"pollvoting/pkg/services"

	"github.com/joho/godotenv"
)

func Serve() {
	godotenv.Load()
	config.SetConfig()

	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = database.AutoMigrate(db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rdb := database.NewRedis(config.LocalConfig.RDBAddress, config.LocalConfig.RDBPassword)
	err = rdb.Ping(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mux := http.NewServeMux()

	//repositories
	userRepo := repositories.NewUserRepository(db, rdb.Client)
	pollRepo := repositories.NewPollRepository(db)

	//service
	userService := services.NewUserService(userRepo)
	pollService := services.NewPollService(pollRepo)

	//controller
	userControllers := controllers.NewUserController(userService)
	pollControllers := controllers.NewPollControllers(pollService)

	routes.Router(mux, userControllers,pollControllers)

	fmt.Println("Server running on port", config.LocalConfig.AppPort)
	log.Fatal(http.ListenAndServe(config.LocalConfig.AppPort, mux))
}
