package Test

import (
	"Bookmarkmanager-Server/Helpers"
	"Bookmarkmanager-Server/Models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http/httptest"
)

var Users = map[string]*Models.User{
	"Administrator": &Models.User{
		Model:         gorm.Model{ID: 1},
		Email:         "admin@test.intern",
		Name:          "Administrator",
		Password:      "admin",
		Administrator: true,
	},
	"User": &Models.User{
		Model:            gorm.Model{ID: 2},
		Email:            "user@test.intern",
		Name:             "User",
		Password:         "user",
		Administrator:    false,
		CategoriesAccess: nil,
	},
	"Jakob": &Models.User{
		Model:            gorm.Model{ID: 3},
		Email:            "jakob@test.intern",
		Name:             "Jakob",
		Password:         "test",
		Administrator:    false,
		CategoriesAccess: nil,
	},
}

var Categories = map[string]*Models.Category{
	"IT": &Models.Category{
		Model:    gorm.Model{ID: 1},
		ParentID: 0,
		Name:     "IT",
		Shared:   false,
		OwnerID:  2,
	},
	"Books": &Models.Category{
		Model:    gorm.Model{ID: 2},
		ParentID: 0,
		Name:     "Books",
		Shared:   true,
		OwnerID:  1,
	},
	"Programming": &Models.Category{
		Model:    gorm.Model{ID: 3},
		ParentID: 1,
		Name:     "Programming",
		Shared:   false,
		OwnerID:  2,
	},
	"C#": &Models.Category{
		Model:    gorm.Model{ID: 4},
		ParentID: 3,
		Name:     "C#",
		Shared:   false,
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
	Models.DatabaseConfig()
	PopulateDatabase()
}

func CreateUsers() {
	var user *Models.User
	user = Users["Administrator"]
	user.Password, _ = Helpers.CreateHashFromPassword(user.Password)
	Models.Database.Create(user)
	user = Users["User"]
	user.Password, _ = Helpers.CreateHashFromPassword(user.Password)
	Models.Database.Create(user)
	user = Users["Jakob"]
	user.Password, _ = Helpers.CreateHashFromPassword(user.Password)
	Models.Database.Create(user)
}

func PopulateDatabase() {

	CreateUsers()

	Models.Database.Create(Categories["IT"])
	Models.Database.Create(Categories["Books"])
	Models.Database.Create(Categories["Programming"])
	Models.Database.Create(Categories["C#"])

	Models.Database.Model(Categories["Books"]).Association("UsersAccess").Append(Users["Administrator"])
	Models.Database.Model(Categories["Books"]).Association("UsersAccess").Append(Users["User"])
	Models.Database.Model(Categories["IT"]).Association("UsersAccess").Append(Users["User"])
	Models.Database.Model(Categories["Programming"]).Association("UsersAccess").Append(Users["User"])
	Models.Database.Model(Categories["C#"]).Association("UsersAccess").Append(Users["User"])
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
	username := "admin@test.intern"
	password := "admin"

	c := GetGinContext(username, password, method, route, body)
	c.Set("UserID", uint(1))

	return c
}

func GetGinContextAsUser(method string, route string, body interface{}) *gin.Context {
	username := "user@test.intern"
	password := "user"

	c := GetGinContext(username, password, method, route, body)
	c.Set("UserID", uint(2))
	return c
}

func UnmarshalObject(c *gin.Context, object interface{}) error {
	err := json.Unmarshal([]byte(c.Writer.(*bodyLogWriter).body.String()), &object)
	return err

}
