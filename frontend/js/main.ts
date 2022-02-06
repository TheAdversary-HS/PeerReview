let articleParent = document.getElementById("articles");

function updateSeach(value: string) {
    if(value == "") {
        addRecent()
        return
    }

    let query: SearchQuery = {
        query: value,
        limit: 5
    }


    clearArticles()
    search(query).then(function(data) {
        data.forEach(function(article) {
            addArticle(article)
        })
    })
}

function clearArticles() {
    articleParent.innerHTML = ""
}

window.onload = function() {
    addRecent()
}

function addRecent () {
    clearArticles()
    recent(5).then(function(data) {
        data.forEach(function(article) {
            addArticle(article)
        })
    })
}

function addArticle(article: ArticleSummary) {
    let articleA = document.createElement("a")
    articleA.setAttribute("href", article.link)

    let articleDiv = document.createElement("div")
    articleDiv.setAttribute("class", "article")

    let articleHeader = document.createElement("div")
    articleHeader.setAttribute("class", "article-header")

    let articleHeaderTitle = document.createElement("h3")
    articleHeaderTitle.innerHTML = article.title

    articleHeader.appendChild(articleHeaderTitle)
    articleDiv.appendChild(articleHeader)

    let articleDescription = document.createElement("div")
    articleDescription.setAttribute("class", "article-description")

    let articleDescriptionP = document.createElement("p")

    let articleDescriptionTopics = document.createElement("i")
    articleDescriptionTopics.innerHTML = article.tags.join(", ")

    let articleDescriptionAuthors = document.createElement("i")
    articleDescriptionAuthors.innerHTML = article.authors[0].name //TODO join ALL Auhtors

    let articleDescriptionDate = document.createElement("i")
    articleDescriptionDate.innerHTML = article.modified.toString()

    articleDescriptionP.appendChild(articleDescriptionTopics)
    articleDescriptionP.appendChild(articleDescriptionAuthors)
    articleDescriptionP.appendChild(articleDescriptionDate)

    articleDescription.appendChild(articleDescriptionP)
    articleDiv.appendChild(articleDescription)

    let articleBody = document.createElement("div")
    articleBody.setAttribute("class", "article-body")

    let articleBodyP = document.createElement("p")
    articleBodyP.innerHTML = article.summary

    articleBody.appendChild(articleBodyP)
    articleDiv.appendChild(articleBody)

    articleA.appendChild(articleDiv)
    articleParent.appendChild(articleA)

    let divider = document.createElement("div")
    divider.setAttribute("class", "divider")

    articleParent.appendChild(divider)
}
