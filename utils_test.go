package main

import (
	"Bookmarkmanager-Server/Handlers"
	"Bookmarkmanager-Server/Models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"gorm.io/gorm"
	"net/http/httptest"
	"testing"
)

var Users = map[string]*Models.User{
	"Administrator": &Models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name:          "Administrator",
		Password:      "admin",
		Administrator: true,
	},
	"User": &Models.User{
		Model: gorm.Model{
			ID: 2,
		},
		Name:             "User",
		Password:         "user",
		Administrator:    false,
		CategoriesAccess: nil,
	},
}

var Categories = map[string]*Models.Category{
	"IT": &Models.Category{
		Model: gorm.Model{
			ID: 1,
		},
		ParentID: 0,
		Name:     "IT",
		Shared:   false,
		OwnerID:  2,
	},
	"Books": &Models.Category{
		Model: gorm.Model{
			ID: 2,
		},
		ParentID: 0,
		Name:     "Books",
		Shared:   true,
		OwnerID:  1,
	},
	"Programming": &Models.Category{
		Model: gorm.Model{
			ID: 3,
		},
		ParentID: 1,
		Name:     "Programming",
		Shared:   true,
		OwnerID:  2,
	},
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func InitializeDatabase() {
	Models.DatabaseConfig(Models.Sqlite, Models.Test)
}

func PopulateDatabase() {
	InitializeDatabase()

	Models.Database.Create(Users["Administrator"])
	Models.Database.Create(Users["User"])

	Models.Database.Create(Categories["IT"])
	Models.Database.Create(Categories["Books"])
	Models.Database.Create(Categories["Programming"])

	Models.Database.Model(Categories["Books"]).Association("UsersAccess").Append(Users["Administrator"])
	Models.Database.Model(Categories["Books"]).Association("UsersAccess").Append(Users["User"])
	Models.Database.Model(Categories["IT"]).Association("UsersAccess").Append(Users["User"])
	Models.Database.Model(Categories["Programming"]).Association("UsersAccess").Append(Users["User"])
}
func CleanupDatabase() {
	Models.Database.Exec("Drop table user_categories")
	Models.Database.Exec("Drop table bookmarks")
	Models.Database.Exec("Drop table categories")
	Models.Database.Exec("Drop table users")
}

func GetGinContext(username string, password string, method string, route string, body interface{}) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	var buf *bytes.Buffer
	if body != nil {
		jsonValue, _ := json.Marshal(body)
		buf = bytes.NewBuffer(jsonValue)
	} else {
		buf = new(bytes.Buffer)
	}

	c.Request = httptest.NewRequest(method, route, buf)
	c.Request.SetBasicAuth(username, password)
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	return c
}

func GetGinContextAsAdministrator(method string, route string, body interface{}) *gin.Context {
	username := "Administrator"
	password := "admin"

	c := GetGinContext(username, password, method, route, body)
	c.Set("UserID", uint(1))

	return c
}

func GetGinContextAsUser(method string, route string, body interface{}) *gin.Context {
	username := "User"
	password := "user"

	c := GetGinContext(username, password, method, route, body)
	c.Set("UserID", uint(2))
	return c
}

func UnmarshalObject(c *gin.Context, object interface{}) error {
	err := json.Unmarshal([]byte(c.Writer.(*bodyLogWriter).body.String()), &object)
	return err

}

func TestGetCurrentUser(t *testing.T) {

	PopulateDatabase()
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
	PopulateDatabase()
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
	PopulateDatabase()
	defer CleanupDatabase()

	shoppingSitesCategory := Models.Category{
		Model:    gorm.Model{ID: 4},
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

func TestEditCategory(t *testing.T) {
	PopulateDatabase()
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
