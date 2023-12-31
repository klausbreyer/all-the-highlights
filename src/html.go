package src

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/klausbreyer/grr"
)

func getHtml(books []Book) template.HTML {
	renderedBooks := make([]template.HTML, len(books))
	for i, book := range books {
		renderedBooks[i] = getBook(DataBook{Book: book, BookIndex: i})
	}

	return grr.Render(`
			<html>
		<head>
			<meta charset="UTF-8">
			<title>all-the-highlights</title>
			{{.Css}}
		</head>
		<body>
				<script>
		function copyToClipboard(copyText) {
			var textArea = document.createElement("textarea");
			textArea.value = copyText;
			document.body.appendChild(textArea);
			textArea.select();
			document.execCommand("Copy");
			textArea.remove();
		}
		</script>
		{{.Head}}
			{{.Books}}
		</body>
		</html>
	`,
		struct {
			Css   template.HTML
			Head  template.HTML
			Books template.HTML
		}{
			getCss(),
			GetHead(),
			grr.Flatten(renderedBooks),
		})
}

func GetHead() template.HTML {
	return grr.Yield(`
	    <h1>
        All the Highlights from
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
            <a href="https://obsidian.md/" target="_blank">Obsidian</a>
            .
            <a href="https://v01.io/2020/12/31/pocket-highlights/" target="_blank">Read more in my blog post</a> (outdated)
            .
        </i>
    </p>
	`)
}

func getCss() template.HTML {
	return grr.Yield(`
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
	// Erstellen Sie die eigene Variable und entfernen Sie alle ":"
	titleStr := strings.ReplaceAll(fmt.Sprintf("%s, %s, %d", data.Book.Title, data.Book.Author, data.Book.FirstHighlightYear), ":", "")

	return grr.Render(`
<a id="{{.BookIndex}}" href="#{{.BookIndex}}">#{{.BookIndex}}</a>
		<h2 style="cursor:copy;" onclick="copyToClipboard('read/{{.TitleStr}}')" >{{.TitleStr}}</h2>
		<span onclick="copyToClipboard('{{.Book.SourceURL}}')" style="cursor:copy;" >{{.Book.SourceURL}}</span>
		<a href="{{.Book.SourceURL}}" target="_blank">&raquo;&raquo;&raquo;</a>
		<ul>
		{{.Highlights}}
		</ul>
		<hr/>
	`, struct {
		Book       Book
		BookIndex  int
		Highlights template.HTML
		TitleStr   string
	}{
		data.Book,
		data.BookIndex,
		getHighlights(data.Book.Highlights),
		titleStr,
	})
}

func getHighlights(highlights []Highlight) template.HTML {
	return grr.Map(`
				<li> {{.Text | html}}</li>

	`, highlights)
}
