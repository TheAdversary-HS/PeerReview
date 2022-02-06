package schema

type ArticleSummary struct {
	Id       int      `json:"id" gorm:"primaryKey"`
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	Authors  []Author `json:"authors" gorm:"-"`
	Image    string   `json:"image"`
	Tags     []string `json:"tags" gorm:"-"`
	Created  int64    `json:"created"`
	Modified int64    `json:"modified"`
	Link     string   `json:"link"`
}

type Author struct {
	Id          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Password    string `json:"-"`
	Information string `json:"information"`
}

type Asset struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}
