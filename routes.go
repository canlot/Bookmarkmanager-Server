package main

import "github.com/canlot/Bookmarkmanager-Server/Handlers"

func initializeRoutes() {
	apiRoutes := router.Group("/apiv1")
	{
		categoryRoutes := apiRoutes.Group("/categories")
		{
			categoryRoutes.GET("/", Handlers.GetAllCategories)
			categoryRoutes.GET("/:id")
			categoryRoutes.POST("/")
			categoryRoutes.PUT("/:id")
			categoryRoutes.DELETE("/:id")
		}

		bookmarkRoutes := apiRoutes.Group("/bookmarks")
		{
			bookmarkRoutes.GET("/")
			bookmarkRoutes.GET("/:id")
			bookmarkRoutes.POST("/")
			bookmarkRoutes.PUT("/:id")
			bookmarkRoutes.DELETE("/:id")
		}

		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.GET("/")
			userRoutes.GET("/:id")
			userRoutes.POST("/")
			userRoutes.PUT("/:id")
			userRoutes.DELETE("/:id")
		}
	}
}
