package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	grr "github.com/klausbreyer/grr"
)

var token string

type Highlight struct {
	ID            int      `json:"id"`
	Text          string   `json:"text"`
	Location      int      `json:"location"`
	LocationType  string   `json:"location_type"`
	Note          string   `json:"note"`
	Color         string   `json:"color"`
	HighlightedAt string   `json:"highlighted_at"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	ExternalID    string   `json:"external_id"`
	EndLocation   string   `json:"end_location"`
	URL           string   `json:"url"`
	BookID        int      `json:"book_id"`
	Tags          []string `json:"tags"`
	IsFavorite    bool     `json:"is_favorite"`
	IsDiscard     bool     `json:"is_discard"`
	ReadwiseURL   string   `json:"readwise_url"`
}

type Book struct {
	UserBookID         int         `json:"user_book_id"`
	Title              string      `json:"title"`
	Author             string      `json:"author"`
	ReadableTitle      string      `json:"readable_title"`
	Source             string      `json:"source"`
	CoverImageURL      string      `json:"cover_image_url"`
	UniqueURL          string      `json:"unique_url"`
	BookTags           []string    `json:"book_tags"`
	Category           string      `json:"category"`
	DocumentNote       string      `json:"document_note"`
	ReadwiseURL        string      `json:"readwise_url"`
	SourceURL          string      `json:"source_url"`
	ASIN               string      `json:"asin"`
	Highlights         []Highlight `json:"highlights"`
	FirstHighlightYear int
}

func fetchFromExportApi(updatedAfter string) ([]Book, error) {
	var fullData []Book
	var nextPageCursor string

	for {
		baseUrl := "https://readwise.io/api/v2/export/"

		params := url.Values{}
		if nextPageCursor != "" {
			params.Add("pageCursor", nextPageCursor)
		}
		if updatedAfter != "" {
			params.Add("updatedAfter", updatedAfter)
		}

		fmt.Println("Making export api request with params " + params.Encode())

		client := &http.Client{}
		req, _ := http.NewRequest("GET", baseUrl+"?"+params.Encode(), nil)
		req.Header.Add("Authorization", "Token "+token)

		response, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		data, _ := ioutil.ReadAll(response.Body)

		var result map[string]interface{}
		json.Unmarshal(data, &result)

		results, ok := result["results"].([]interface{})
		fmt.Println("Got", len(results), "results")

		if ok && len(results) > 0 {
			for _, item := range results {
				var book Book
				jsonData, _ := json.Marshal(item)
				json.Unmarshal(jsonData, &book)
				fullData = append(fullData, book)
			}
		}

		nextPageCursor, ok = result["nextPageCursor"].(string)
		//@todo. check why it is not paging.
		if !ok || nextPageCursor == "" {
			break
		}
	}

	return fullData, nil
}

func getHtml(books []Book, yield ...template.HTML) template.HTML {
	renderedBooks := make([]template.HTML, len(books))
	for i, book := range books {
		renderedBooks[i] = getBook(DataBook{Book: book, BookIndex: i})
	}

	allRenderedBooks := grr.Flatten(renderedBooks)

	dm := grr.Extend(nil)
	dm["Css"] = getCss()
	dm["Head"] = GetHead()
	return grr.Render(`
			<html>
		<head>
			<meta charset="UTF-8">
			<title>all-the-highlights</title>
			{{.Css}}
		</head>
		<body>
		{{.Head}}
			{{yield}}
		</body>
		</html>
	`, dm, allRenderedBooks)
}

type DataBook struct {
	Book      Book
	BookIndex int
}

func GetHead() template.HTML {
	return grr.Render(`
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
	`, nil)
}

func getCss() template.HTML {
	return grr.Render(`
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
	`, nil)
}

func getBook(data DataBook, yield ...template.HTML) template.HTML {
	highlights := getHighlights(DataHighlight{Highlights: data.Book.Highlights})

	// dataMap := grr.Extend(data)
	// dataMap["Foot"] = getFoot(DataFoot{Copy: "© 2021"})
	return grr.Render(`
		<h1> {{.Book.Title}}, {{.Book.Author}}, {{.Book.FirstHighlightYear}}</h1>
			<a href="#{{.BookIndex}}">#{{.BookIndex}}</a>
			{{yield}}
		<hr/>
	`, data, highlights)
}

type DataHighlight struct {
	Highlights []Highlight
}

func getHighlights(data DataHighlight) template.HTML {
	return grr.Render(`
		<ul>
				{{range $highlightIndex, $highlight := .Highlights}}
					<li> {{$highlight.Text | html}}</li>
				{{end}}
				</ul>
	`, data)
}

func toHTML(books []Book) error {
	html := getHtml(books)
	file, err := os.Create("dist/index.html")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(html))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	token = os.Getenv("READWISE_TOKEN")

	allData, err := fetchFromExportApi("")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Filter out books without highlights
	var booksWithHighlights []Book
	for _, book := range allData {
		if len(book.Highlights) > 0 {
			firstHighlightDate, _ := time.Parse(time.RFC3339, book.Highlights[0].HighlightedAt)
			book.FirstHighlightYear = firstHighlightDate.Year()
			booksWithHighlights = append(booksWithHighlights, book)
		}
	}

	// Sort books by the date of the last highlight
	sort.Slice(booksWithHighlights, func(i, j int) bool {
		lastHighlightI, _ := time.Parse(time.RFC3339, booksWithHighlights[i].Highlights[len(booksWithHighlights[i].Highlights)-1].HighlightedAt)
		lastHighlightJ, _ := time.Parse(time.RFC3339, booksWithHighlights[j].Highlights[len(booksWithHighlights[j].Highlights)-1].HighlightedAt)
		return lastHighlightI.Before(lastHighlightJ)
	})

	err = toHTML(booksWithHighlights)
	if err != nil {
		fmt.Println("Error creating HTML:", err)
	}
}
