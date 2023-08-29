package main

import "github.com/canlot/Bookmarkmanager-Server/Handlers"

func InitializeRoutes() {
	apiRoutes := router.Group("/apiv1", Handlers.Authenticate)
	{
		apiRoutes.GET("/currentuser", Handlers.GetCurrentUser)

		categoryRoutes := apiRoutes.Group("/categories")
		{
			categoryRoutes.GET("/", Handlers.GetCategories)
			categoryRoutes.POST("/", Handlers.AddCategory)

			categoryRoutesID := categoryRoutes.Group("/:category_id")
			{
				categoryRoutesID.GET("/", Handlers.GetCategories)
				categoryRoutesID.PUT("/", Handlers.EditCategory)
				categoryRoutesID.DELETE("/", Handlers.DeleteCategory)

				categoryRoutesIDBookmarks := categoryRoutesID.Group("/bookmarks")
				{
					categoryRoutesIDBookmarks.GET("/", Handlers.GetBookmarksWithCategoryId)
					categoryRoutesIDBookmarks.POST("/", Handlers.AddBookmarkToCategory)
				}
				categoryRoutesIDUsers := categoryRoutesID.Group("/permissions")
				{
					categoryRoutesIDUsers.GET("/", Handlers.GetUsersForCategory)
					categoryRoutesIDUsers.POST("/", Handlers.AddUsersForCategoryOnce)
					categoryRoutesIDUsers.POST("/inherit", Handlers.AddUsersForCategoryInherit)
					categoryRoutesIDUsers.DELETE("/", Handlers.RemoveUsersFromCategory)
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
			userRoutes.GET("/", Handlers.GetAllUsers)
			userRoutes.GET("/:id")
			userRoutes.POST("/")
			userRoutes.PUT("/:id")
			userRoutes.DELETE("/:id")
		}
	}
}
