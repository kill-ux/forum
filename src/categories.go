package forum

import (
	"log"
	"strings"
)

// FillCategoriesDB populates the categories table in the database with default values.
func FillCategoriesDB() {
	query := `
		INSERT INTO categories (name)
		VALUES 
			("art"),
			("programming"),
			("news"),
			("studying"),
			("business"),
			("Discussions"),
			("Questions"),
			("Ideas"),
			("Articles"),
			("Events"),
			("Issues"),
			("Others");
	`

	// Execute the query and handle errors
	_, err := DB.Exec(query)
	if err != nil {
		// Ignore "UNIQUE constraint failed" errors to avoid duplicate inserts
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			log.Println("Categories already exist, skipping insertion.")
		} else {
			panic("Error inserting categories: " + err.Error())
		}
	}
}

// FillCategories fetches all categories from the database and populates the Page struct.
func (data *Page) FillCategories() {
	query := "SELECT id, name FROM categories"

	// Execute the query
	rows, err := DB.Query(query)
	if err != nil {
		Mux.Lock()
		DB.Close()
		Mux.Unlock()
		log.Fatalf("Error fetching categories: %v", err)
	}
	defer rows.Close()

	// Iterate over the rows and populate the Page struct
	for rows.Next() {
		var category Categorie
		if err := rows.Scan(&category.Id, &category.CatName); err != nil {
			log.Printf("Error scanning category row: %v", err)
			continue
		}
		data.Categories = append(data.Categories, category)
	}
}
