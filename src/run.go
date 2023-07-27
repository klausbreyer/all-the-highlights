package src

import (
	"fmt"
	"os"
	"sort"
	"time"
)

var token string

func writeHtml(books []Book) error {
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

func Run() {
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

	err = writeHtml(booksWithHighlights)
	if err != nil {
		fmt.Println("Error creating HTML:", err)
	}
}
