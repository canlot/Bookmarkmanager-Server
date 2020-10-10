package Models

type Category struct {
	ID     uint   `json:"id"`
	TopID  uint   `json:"topid"`
	Name   string `json:"name"`
	Shared bool   `json:"shared"`
}

var categoryList = []Category{
	{
		ID:     1,
		TopID:  0,
		Name:   "IT",
		Shared: false,
	},
	{
		ID:     2,
		TopID:  1,
		Name:   "Programmiersprachen",
		Shared: false,
	},
}

func GetAllCategories() []Category {
	return categoryList
}
