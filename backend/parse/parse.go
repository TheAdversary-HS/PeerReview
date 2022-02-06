package parse

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func newParser() *parser.Parser {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs

	return parser.NewWithExtensions(extensions)
}

func newHtmlRenderer() *html.Renderer {
	renderOpts := html.RendererOptions{
		Flags: html.CommonFlags | html.LazyLoadImages,
	}
	return html.NewRenderer(renderOpts)
}

func ParseToHtml(rawMarkdown []byte) []byte {
	node := markdown.Parse(rawMarkdown, newParser())
	return markdown.Render(node, newHtmlRenderer())
}
