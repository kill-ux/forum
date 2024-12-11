package forum

import (
	"log"
)

func (data *Page) FillCategories() {
	query := "SELECT * FROM categories"
	rows, err := data.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		category := Categorie{}
		rows.Scan(&category.Id, &category.CatName)
		data.Categories = append(data.Categories, category)
	}
}
