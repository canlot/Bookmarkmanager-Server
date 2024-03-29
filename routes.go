package main

import "Bookmarkmanager-Server/Handlers"

func InitializeRoutes() {
	router.POST("/apiv1/login", Handlers.GetBearerToken)

	apiRoutes := router.Group("/apiv1", Handlers.Authenticate)
	{
		apiRoutes.GET("/currentuser", Handlers.GetCurrentUser)

		categoryRoutes := apiRoutes.Group("/categories")
		{
			categoryRoutes.GET("/", Handlers.GetCategories)
			categoryRoutes.POST("/", Handlers.AddCategory)
			categoryRoutes.GET("/search/:search_text", Handlers.SearchCategories)

			categoryRoutesID := categoryRoutes.Group("/:category_id")
			{
				categoryRoutesID.GET("/", Handlers.GetCategories)
				categoryRoutesID.PUT("/", Handlers.EditCategory)
				categoryRoutesID.DELETE("/", Handlers.DeleteCategory)

				categoryRoutesIDBookmarks := categoryRoutesID.Group("/bookmarks")
				{
					categoryRoutesIDBookmarks.GET("/", Handlers.GetBookmarksWithCategoryId)
					categoryRoutesIDBookmarks.POST("/", Handlers.AddBookmarkToCategory)
					categoryRoutesIDBookmarks.PUT("/:bookmark_id", Handlers.EditBookmarkWithBookmarkId)
					categoryRoutesIDBookmarks.DELETE("/:bookmark_id", Handlers.DeleteBookmarkWithBookmarkId)
				}
				categoryRoutesIDUsers := categoryRoutesID.Group("/permissions")
				{
					categoryRoutesIDUsers.GET("/", Handlers.GetUsersForCategory)
					categoryRoutesIDUsers.PUT("/", Handlers.EditUsersForCategory)
				}
			}
		}

		bookmarkRoutes := apiRoutes.Group("/bookmarks")
		{
			bookmarkRoutes.GET("/:id")
			bookmarkRoutes.PUT("/:id")
			bookmarkRoutes.DELETE("/:id")

			bookmarkRoutes.GET("/search/:search_text", Handlers.SearchBookmarks)
		}

		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.GET("/", Handlers.GetAllUsers)
			userRoutes.GET("/search/:search_text", Handlers.SearchUsers)
			userRoutes.GET("/:id")
			userRoutes.POST("/:password", Handlers.AddUser)
			userRoutes.PUT("/:id", Handlers.EditUser)
			userRoutes.PUT("/:id/:password", Handlers.EditUser)
			userRoutes.DELETE("/:id", Handlers.DeleteUser)
		}
	}
}
