package models

type Book struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Author    *Author `json:"author"`
	Year      int     `json:"year"`
	Available bool    `json:"available"`
}
