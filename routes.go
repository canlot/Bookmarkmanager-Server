package main

import "github.com/canlot/Bookmarkmanager-Server/Handlers"

func initializeRoutes() {
	apiRoutes := router.Group("/apiv1", Handlers.Authenticate)
	{
		categoryRoutes := apiRoutes.Group("/categories")
		{
			categoryRoutes.GET("/", Handlers.GetCategories)
			categoryRoutes.POST("/", Handlers.AddCategory)

			categoryRoutesID := categoryRoutes.Group("/:category_id")
			{
				categoryRoutesID.GET("/", Handlers.GetCategories)
				categoryRoutesID.PUT("/")
				categoryRoutesID.DELETE("/")

				categoryRoutesIDBookmarks := categoryRoutesID.Group("/bookmarks")
				{
					categoryRoutesIDBookmarks.GET("/", Handlers.GetBookmarksWithCategoryId)
					categoryRoutesIDBookmarks.POST("/", Handlers.AddBookmarkToCategory)
				}
				categoryRoutesIDUsers := categoryRoutesID.Group("/users")
				{
					categoryRoutesIDUsers.GET("/", Handlers.GetUsersForCategoryFull)
					categoryRoutesIDUsers.GET("/inherit", Handlers.GetUsersForCategoryInherit)
					categoryRoutesIDUsers.POST("/", Handlers.AddUsersForCategory)
					categoryRoutesIDUsers.DELETE("/:user_id")
				}
			}
		}

		bookmarkRoutes := apiRoutes.Group("/bookmarks")
		{
			bookmarkRoutes.GET("/")
			bookmarkRoutes.GET("/:id")
			bookmarkRoutes.POST("/:id")
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
