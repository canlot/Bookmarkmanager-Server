package Test

import (
	"Bookmarkmanager-Server/Handlers"
	"Bookmarkmanager-Server/Models"
	"github.com/go-playground/assert/v2"
	"gorm.io/gorm"
	"testing"
)

func TestGetCurrentUser(t *testing.T) {

	InitializeDatabase()
	defer CleanupDatabase()

	c := GetGinContextAsUser("GET", "/currentuser", nil)

	Handlers.GetCurrentUser(c)

	assert.Equal(t, 200, c.Writer.Status())
	assert.Equal(t, c.Writer.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var user Models.User
	if err := UnmarshalObject(c, &user); err != nil {
		t.Error(err)
	}
	assert.Equal(t, user.Name, Users["User"].Name)
	assert.Equal(t, user.Password, "")
}

func TestGetCategory(t *testing.T) {
	InitializeDatabase()
	defer CleanupDatabase()

	c := GetGinContextAsUser("GET", "/apiv1/categories", nil)
	//c.AddParam("category_id", "1")
	Handlers.GetCategories(c)

	assert.Equal(t, 200, c.Writer.Status())
	assert.Equal(t, c.Writer.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var category []Models.Category
	if err := UnmarshalObject(c, &category); err != nil {
		t.Error(err)
	}

	assert.Equal(t, category[0].OwnerID, Categories["IT"].OwnerID)
	assert.Equal(t, category[1].OwnerID, Categories["Books"].OwnerID)
	assert.Equal(t, category[0].Name, Categories["IT"].Name)
	assert.Equal(t, category[1].Name, Categories["Books"].Name)
}

func TestAddCategory(t *testing.T) {
	InitializeDatabase()
	defer CleanupDatabase()

	shoppingSitesCategory := Models.Category{
		Model:    gorm.Model{ID: 5},
		ParentID: 0,
		Name:     "Shopping",
		Shared:   false,
		OwnerID:  1,
	}

	c := GetGinContextAsUser("POST", "/apiv1/categories", &shoppingSitesCategory)
	Handlers.AddCategory(c)

	assert.Equal(t, 200, c.Writer.Status())

	var category Models.Category

	Models.Database.Take(&category, shoppingSitesCategory.ID)

	assert.Equal(t, category.Name, shoppingSitesCategory.Name)
	assert.Equal(t, category.OwnerID, Users["User"].ID) // Should be set to "Users" "Id" because "User" created it

}

func TestAddCategoryChild(t *testing.T) {
	InitializeDatabase()
	defer CleanupDatabase()

	golangCategory := Models.Category{
		Model:    gorm.Model{ID: 5},
		ParentID: 3,
		Name:     "Golang",
		Shared:   false,
		OwnerID:  1, //wrong OwnerId at start, should be 2 -> User
	}

	//test with wrong OwnerId
	c := GetGinContextAsAdministrator("POST", "/apiv1/categories", &golangCategory)
	Handlers.AddCategory(c)

	assert.Equal(t, 400, c.Writer.Status())

	var testCategory Models.Category
	if dbctx := Models.Database.Take(&testCategory, golangCategory.ID); dbctx.Error == nil {
		t.Error("Category has been still added")
	}

	//test with right OwnerId
	golangCategory.OwnerID = 2
	c = GetGinContextAsUser("POST", "/apiv1/categories", &golangCategory)
	Handlers.AddCategory(c)

	assert.Equal(t, 200, c.Writer.Status())

	var category Models.Category

	Models.Database.Take(&category, golangCategory.ID)

	assert.Equal(t, category.Name, golangCategory.Name)
	assert.Equal(t, category.OwnerID, Users["User"].ID) // Should be set to "Users" "Id" because "User" created it

}

func TestEditCategory(t *testing.T) {
	InitializeDatabase()
	defer CleanupDatabase()
	testCategory := Models.Category{
		Model: gorm.Model{
			ID: 2,
		},
		ParentID: 0,
		Name:     "BooksTest",
		Shared:   true,
		OwnerID:  1,
	}

	//valid case:
	c := GetGinContextAsAdministrator("PUT", "/apiv1/categories", &testCategory)
	Handlers.EditCategory(c)

	assert.Equal(t, 200, c.Writer.Status())

	var category Models.Category

	Models.Database.Take(&category, testCategory.ID)

	assert.Equal(t, category.Name, testCategory.Name)
	assert.Equal(t, category.OwnerID, Users["Administrator"].ID) // Should be set to "Users" "Id" because "User" created it

	//case 1: user is not the owner of the category

	c = GetGinContextAsUser("PUT", "/apiv1/categories", &testCategory)
	Handlers.EditCategory(c)

	assert.Equal(t, 400, c.Writer.Status())

	//case 2: user is owner but the OwnerID is not right

	testCategory.OwnerID = Users["User"].ID

	c = GetGinContextAsAdministrator("PUT", "/apiv1/categories", &testCategory)
	Handlers.EditCategory(c)

	assert.Equal(t, 400, c.Writer.Status())

	//case 3: ParentID changed

	testCategory.OwnerID = Categories["Books"].OwnerID
	testCategory.ParentID = 1

	c = GetGinContextAsAdministrator("PUT", "/apiv1/categories", &testCategory)
	Handlers.EditCategory(c)

	assert.Equal(t, 400, c.Writer.Status())

	//case 4: not existent id
	testCategory.ID = 100
	testCategory.ParentID = Categories["Books"].ParentID

	c = GetGinContextAsAdministrator("PUT", "/apiv1/categories", &testCategory)
	Handlers.EditCategory(c)

	assert.Equal(t, 400, c.Writer.Status())

}

func TestDeleteCategory(t *testing.T) {
	InitializeDatabase()
	defer CleanupDatabase()
	//case 1: wrong user, that has no access on category
	c := GetGinContextAsAdministrator("DELETE", "/apiv1/categories", nil)
	c.AddParam("category_id", "3")
	Handlers.DeleteCategory(c)

	assert.Equal(t, 400, c.Writer.Status())

	c = GetGinContextAsUser("DELETE", "/apiv1/categories", nil)
	c.AddParam("category_id", "3")
	Handlers.DeleteCategory(c)

	assert.Equal(t, 200, c.Writer.Status())
	var category Models.Category
	if result := Models.Database.Take(&category, 3); result.Error == nil {
		t.Error("Category has not been deleted")
	}
	//TODO: also check if the associations has been deleted
	//Models.Database.Select("")

}

func userExistInList(username string, users []Models.User) bool {
	for i := 0; i < len(users); i++ {
		if users[i].Name == username {
			return true
		}
	}
	return false
}

func TestChangeUserPermissionCategoryWithExistingUsers(t *testing.T) {
	InitializeDatabase()
	defer CleanupDatabase()

	usersWithAccess := []Models.User{
		Models.User{
			Model:         gorm.Model{ID: 1},
			Name:          "Administrator",
			Password:      "admin",
			Administrator: true,
		},
		Models.User{
			Model:         gorm.Model{ID: 3},
			Name:          "Jakob",
			Password:      "test",
			Administrator: false,
		},
	}

	c := GetGinContextAsAdministrator("PUT", "/apiv1/categories/2/permissions", &usersWithAccess)
	c.AddParam("category_id", "2")
	Handlers.EditUsersForCategory(c)

	assert.Equal(t, 200, c.Writer.Status())

	categoryBooks := Models.Category{
		Model:    gorm.Model{ID: 2},
		ParentID: 0,
		Name:     "Books",
		Shared:   true,
		OwnerID:  1,
	}

	var users []Models.User

	Models.Database.Model(&categoryBooks).Association("UsersAccess").Find(&users)

	assert.Equal(t, true, userExistInList("Administrator", users))
	assert.Equal(t, false, userExistInList("User", users))
	assert.Equal(t, true, userExistInList("Jakob", users))

}

func TestChangeUserPermissionCategoryWithNonExistingUsers(t *testing.T) {
	InitializeDatabase()
	defer CleanupDatabase()

	usersWithAccess := []Models.User{
		Models.User{
			Model:         gorm.Model{ID: 1},
			Name:          "Administrator",
			Password:      "admin",
			Administrator: true,
		},
		Models.User{
			Model:         gorm.Model{ID: 3},
			Name:          "Jakob",
			Password:      "test",
			Administrator: false,
		},
		Models.User{
			Model:         gorm.Model{ID: 4},
			Name:          "Bob",
			Administrator: false,
		},
	}

	c := GetGinContextAsAdministrator("PUT", "/apiv1/categories/2/permissions", &usersWithAccess)
	c.AddParam("category_id", "2")
	Handlers.EditUsersForCategory(c)

	categoryBooks := Models.Category{
		Model:    gorm.Model{ID: 2},
		ParentID: 0,
		Name:     "Books",
		Shared:   true,
		OwnerID:  1,
	}

	var users []Models.User

	Models.Database.Model(&categoryBooks).Association("UsersAccess").Find(&users)

	assert.Equal(t, true, userExistInList("Administrator", users))
	assert.Equal(t, true, userExistInList("User", users))
	assert.Equal(t, false, userExistInList("Jakob", users))
	assert.Equal(t, false, userExistInList("Bob", users))

	assert.Equal(t, 400, c.Writer.Status())

}
