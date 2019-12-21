package model

type Cat struct {
	ID    string `json:"id, omitempty"`
	Name  string `json:"name, omitempty"`
	Color string `json:"color, omitempty"`
	Age   int    `json:"age, omitempty"`
}

// GetCat retrieves a single cat from the database
//func (c Cat) GetCat(s Storage) ([]byte, error) {
//	return s.Select(c.ID)
//}
//
//// CreateCat creates a new cat in the database
//func (c Cat) CreateCat(s Storage) (string, error) {
//
//	return s.Insert(c)
//}
//
//// UpdatesCat creates a new cat in the database
//func (c Cat) UpdateCat(s Storage) error {
//	return s.Update(c.ID, c)
//}
//
//// DeleteCat removes a cat in the database
//func (c Cat) DeleteCat(s Storage) error {
//	return s.Delete(c.ID)
//}
//

// GetCats retreives multiple cats from the database
// todo finish implementing this
//func GetCats(s Storage, start, count int) ([]Cat, error) {
//
//	rows, err := db.Query("SELECT id, name, kind,color,age FROM cats LIMIT  $1 OFFSET $2", count, start)
//
//	//if err == nil {
//	//
//	//}
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	cats := []Cat{}
//
//	for rows.Next() {
//		var c Cat
//		if err := rows.Scan(&c.ID, &c.Name, &c.Color, &c.Age); err != nil {
//			return nil, err
//		}
//		cats = append(cats, p)
//	}
//
//	return cats, nil
//
//	////////
//	return s.GetAll()
//}
