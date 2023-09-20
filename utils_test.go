package main

import (
	"Bookmarkmanager-Server/Handlers"
	"Bookmarkmanager-Server/Models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http/httptest"
	"testing"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func InitializeDatabase() {
	Models.DatabaseConfig(Models.Mysql, Models.Test)
}

func PopulateDatabase() {
	InitializeDatabase()
	admin := Models.User{
		Name:             "Administrator",
		Password:         "admin",
		Administrator:    true,
		CategoriesAccess: nil,
	}
	user := Models.User{
		Name:             "User",
		Password:         "user",
		Administrator:    false,
		CategoriesAccess: nil,
	}

	itcategory := Models.Category{
		ParentID: 0,
		Name:     "IT",
		Shared:   false,
		OwnerID:  2,
	}

	Models.Database.Create(&admin)
	Models.Database.Create(&user)

	Models.Database.Create(&itcategory)
	Models.Database.Model(&itcategory).Association("UsersAccess").Append(&user)

}
func CleanupDatabase() {
	Models.Database.Exec("Drop table user_categories")
	Models.Database.Exec("Drop table bookmarks")
	Models.Database.Exec("Drop table categories")
	Models.Database.Exec("Drop table users")
}

func GetGinContext(method string, route string) *gin.Context {
	username := "Administrator"
	password := "admin"
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	buf := new(bytes.Buffer)
	c.Request = httptest.NewRequest(method, route, buf)
	c.Request.SetBasicAuth(username, password)
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	return c
}

func UnmarshalObject(c *gin.Context, object interface{}) error {
	err := json.Unmarshal([]byte(c.Writer.(*bodyLogWriter).body.String()), &object)
	return err

}

func TestGetCurrentUser(t *testing.T) {

	PopulateDatabase()
	defer CleanupDatabase()

	c := GetGinContext("GET", "/currentuser")

	Handlers.GetCurrentUser(c)

	assert.Equal(t, 200, c.Writer.Status())
	assert.Equal(t, c.Writer.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var user Models.User
	if err := UnmarshalObject(c, &user); err != nil {
		t.Error(err)
	}
	assert.Equal(t, user.Name, "Administrator")
	assert.Equal(t, user.Password, "")
}

func TestGetCategory(t *testing.T) {
	PopulateDatabase()
	defer CleanupDatabase()

	c := GetGinContext("GET", "/apiv1/categories/1")
	//c.AddParam("category_id", "1")
	c.Set("UserID", uint(2))
	Handlers.GetCategories(c)

	assert.Equal(t, 200, c.Writer.Status())
	assert.Equal(t, c.Writer.Header().Get("Content-Type"), "application/json; charset=utf-8")

	var category []Models.Category
	if err := UnmarshalObject(c, &category); err != nil {
		t.Error(err)
	}

	assert.Equal(t, category[0].OwnerID, uint(2))
}
