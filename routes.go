package main

import "Bookmarkmanager-Server/Handlers"

func InitializeRoutes() {
	router.POST("/apiv1/login", Handlers.GetBearerToken)

	apiRoutes := router.Group("/apiv1", Handlers.Authenticate)
	{
		apiRoutes.POST("/upload", Handlers.UploadTest)
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
					categoryRoutesIDBookmarks.GET("/:bookmark_id/icon", Handlers.GetIconForBookmark)

					categoryRoutesIDBookmarks.POST("/", Handlers.AddBookmarkToCategory)

					categoryRoutesIDBookmarks.POST("/:bookmark_id/icon", Handlers.UploadIconToBookmark)

					categoryRoutesIDBookmarks.PUT("/:bookmark_id", Handlers.EditBookmarkWithBookmarkId)
					categoryRoutesIDBookmarks.PUT("/:bookmark_id/to/:category_destination_id", Handlers.MoveBookmarkWithBookmarkId)
					categoryRoutesIDBookmarks.DELETE("/:bookmark_id", Handlers.DeleteBookmarkWithBookmarkId)
				}
				categoryRoutesIDUsers := categoryRoutesID.Group("/permissions")
				{
					categoryRoutesIDUsers.GET("/", Handlers.GetUsersForCategory)
					categoryRoutesIDUsers.PUT("/", Handlers.EditUsersForCategory)
				}
			}
		}

		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.GET("/", Handlers.GetAllUsers)
			userRoutes.GET("/search/:search_text", Handlers.SearchUsers)
			userRoutes.GET("/:id")
			userRoutes.POST("/:password", Handlers.AddUser)
			userRoutes.PUT("/:id", Handlers.EditUser)
			userRoutes.PUT("/:id/:password", Handlers.SetPassword)
			userRoutes.DELETE("/:id", Handlers.DeleteUser)
		}
	}
}
