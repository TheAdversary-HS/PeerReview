package database

type Article struct {
	Id       int `gorm:"primaryKey"`
	Title    string
	Summary  string
	Image    string
	Created  int64
	Modified int64
	Link     string
	Markdown string
	Html     string
}
