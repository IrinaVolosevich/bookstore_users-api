package app

import "../controllers/ping"
import "../controllers/users"

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
	//router.GET("/users/search", controllers.SearchUser)
}