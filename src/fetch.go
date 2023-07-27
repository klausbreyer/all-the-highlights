package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Highlight struct {
	ID            int    `json:"id"`
	Text          string `json:"text"`
	Location      int    `json:"location"`
	LocationType  string `json:"location_type"`
	Note          string `json:"note"`
	Color         string `json:"color"`
	HighlightedAt string `json:"highlighted_at"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ExternalID    string `json:"external_id"`
	EndLocation   string `json:"end_location"`
	URL           string `json:"url"`
	BookID        int    `json:"book_id"`
	Tags          []Tag  `json:"tags"`
	IsFavorite    bool   `json:"is_favorite"`
	IsDiscard     bool   `json:"is_discard"`
	ReadwiseURL   string `json:"readwise_url"`
}

type BookTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	UserBookID         int         `json:"user_book_id"`
	Title              string      `json:"title"`
	Author             string      `json:"author"`
	ReadableTitle      string      `json:"readable_title"`
	Source             string      `json:"source"`
	CoverImageURL      string      `json:"cover_image_url"`
	UniqueURL          string      `json:"unique_url"`
	BookTags           []BookTag   `json:"book_tags"`
	Category           string      `json:"category"`
	DocumentNote       string      `json:"document_note"`
	ReadwiseURL        string      `json:"readwise_url"`
	SourceURL          string      `json:"source_url"`
	ASIN               string      `json:"asin"`
	Highlights         []Highlight `json:"highlights"`
	FirstHighlightYear int
}

func main() {
	data, err := fetchFromExportApi("")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Received data: ", data)
}

func fetchFromExportApi(updatedAfter string) ([]Book, error) {
	var fullData []Book
	var nextPageCursor string

	fmt.Println("Token: " + token) // Token ausgeben

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
			fmt.Println("Request error: ", err) // Request-Fehler ausgeben
			return nil, err
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Read response error: ", err) // Fehler beim Lesen der Antwort ausgeben
			return nil, err
		}

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			fmt.Println("JSON unmarshal error: ", err) // Fehler beim JSON Unmarshal ausgeben
			return nil, err
		}

		results, ok := result["results"].([]interface{})
		fmt.Println("Got", len(results), "results")

		if ok && len(results) > 0 {
			for _, item := range results {

				var book Book
				jsonData, _ := json.Marshal(item)
				err = json.Unmarshal(jsonData, &book)
				if err != nil {
					fmt.Println("JSON unmarshal book error: ", err) // Fehler beim JSON Unmarshal des Buches ausgeben
					return nil, err
				}
				fullData = append(fullData, book)
			}
		}
		nextPageCursorFloat, ok := result["nextPageCursor"].(float64)
		if ok {
			nextPageCursor = fmt.Sprintf("%.0f", nextPageCursorFloat)
		} else {
			nextPageCursor = ""
		}
		fmt.Println("Next page cursor: ", nextPageCursor) // nextPageCursor ausgeben
		if !ok || nextPageCursor == "" {
			break
		}
	}

	return fullData, nil
}
