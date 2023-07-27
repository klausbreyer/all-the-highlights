package src

import (
	"html/template"

	"github.com/klausbreyer/grr"
)

func getHtml(books []Book) template.HTML {
	renderedBooks := make([]template.HTML, len(books))
	for i, book := range books {
		renderedBooks[i] = getBook(DataBook{Book: book, BookIndex: i})
	}

	return grr.Render(
		struct {
			Css   template.HTML
			Head  template.HTML
			Books template.HTML
		}{
			getCss(),
			GetHead(),
			grr.Flatten(renderedBooks),
		},
		`
			<html>
		<head>
			<meta charset="UTF-8">
			<title>all-the-highlights</title>
			{{.Css}}
		</head>
		<body>
		{{.Head}}
			{{.Books}}
		</body>
		</html>
	`)
}

func GetHead() template.HTML {
	return grr.Render(nil, `
	    <h1>
        Pocket Highlights from
        <a href="https://v01.io" target="_blank">Klaus
                    Breyer</a>
    </h1>
    <p>
        <i>
            This page is created by an personal open source project of mine called

            <a href="https://github.com/klausbreyer/all-the-highlights" target="_blank">all-the-highlights</a>
             to
                        extracts my
                        Highlights from
            <a href="https://read.readwise.io" target="_blank">Readwise Reader</a>
             and format them for
                        easy copy & pastable into
            <a href="https://roamresearch.com" target="_blank">Roam Research</a>
            .
            <a href="https://v01.io/2020/12/31/pocket-highlights/" target="_blank">Read more in my blog post</a>
            .
        </i>
    </p>
	`)
}

func getCss() template.HTML {
	return grr.Render(nil, `
    <style>
    body {
        max-width: 1024px;
        margin: auto;
        font-family: Iowan Old Style, Apple Garamond, Baskerville, Times New Roman, Droid Serif, Times, Source Serif Pro, serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol;
    }

    aside, li {
        line-height: 2;
        font-weight: 600;
        letter-spacing: -0.2px;
    }
    </style>
	`)
}

type DataBook struct {
	Book      Book
	BookIndex int
}

func getBook(data DataBook) template.HTML {
	return grr.Render(struct {
		Book       Book
		BookIndex  int
		Highlights template.HTML
	}{
		data.Book,
		data.BookIndex,
		getHighlights(data.Book.Highlights),
	}, `
		<h1> {{.Book.Title}}, {{.Book.Author}}, {{.Book.FirstHighlightYear}}</h1>
			<a id="{{.BookIndex}}" href="#{{.BookIndex}}">#{{.BookIndex}}</a>
			{{.Highlights}}
		<hr/>
	`)
}

func getHighlights(highlights []Highlight) template.HTML {
	return grr.Render(struct {
		Highlights []Highlight
	}{
		highlights,
	}, `
		<ul>
				{{range $highlightIndex, $highlight := .Highlights}}
					<li> {{$highlight.Text | html}}</li>
				{{end}}
				</ul>
	`)
}
