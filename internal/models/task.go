package models

type Task struct {
	// ID      int64  `json:"id" db:"id"`
	ID      string `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Date    string `json:"date" db:"date"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}
