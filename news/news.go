package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Http     *http.Client
	Key      string
	PageSize int
}

func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}

func (c *Client) FetchArticles(query, page string) (*Results, error) {
	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(query), c.PageSize, page, c.Key)
	response, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := &Results{}
	return res, json.Unmarshal(body, res)
}

type Article struct {
	Source struct {
		ID   int    "json:id"
		Name string "json:name"
	} "json:source"

	Author      string    "json:author"
	Title       string    "json:title"
	Description string    "json:description"
	URL         string    "json:url"
	URLToImage  string    "json:urlToImage"
	PublishedAt time.Time "json:publishedAt"
	Content     string    "json:content"
}

func (a *Article) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}

type Results struct {
	Status       string    "json:status"
	TotalResults int       "json:totalResults"
	Articles     []Article "json:articles"
}
