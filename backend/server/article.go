package server

import (
	"TheAdversary/config"
	"TheAdversary/database"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"path/filepath"
	"text/template"
	"time"
)

var tmpl = template.Must(template.ParseFiles(filepath.Join(config.FrontendDir, "html", "article.gohtml")))

type tmplArticle struct {
	Title    string
	BasePath string
	Summary  string
	Image    string
	Authors  []string
	Tags     []string
	Date     string
	Modified bool
	Content  string
}

func Article(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	articleName := mux.Vars(r)["article"]
	var article database.Article
	if database.GetDB().Table("article").Where("link = ?", articleName).First(&article).RowsAffected == 0 {
		Error404(w, r)
	} else if database.GetDB().Error != nil {
		Error500(w, r)
	} else {
		var authors, tags []string
		database.GetDB().Table("author").Select("name").Where("id IN (?)", database.GetDB().Table("article_author").Select("author_id").Where("article_id = ?", article.Id)).Find(&authors)
		database.GetDB().Table("article_tag").Where("article_id = ?", article.Id).Find(&tags)

		ta := tmplArticle{
			Title:    article.Title,
			BasePath: config.Address + path.Join("/", config.SubPath) + "/",
			Summary:  article.Summary,
			Image:    article.Image,
			Authors:  authors,
			Tags:     tags,
			Content:  article.Html,
		}
		if article.Modified > 0 {
			ta.Date = time.Unix(article.Modified, 0).Format("Monday, 2. January 2006 | 15:04")
			ta.Modified = true
		} else {
			ta.Date = time.Unix(article.Created, 0).Format("Monday, 2. January 2006 | 15:04")
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, ta)
	}
}
