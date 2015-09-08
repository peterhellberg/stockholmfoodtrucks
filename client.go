package stockholmfoodtrucks

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// FoodTruck contains the information for a food truck
type FoodTruck struct {
	Name     string    `json:"name"`
	Text     string    `json:"text"`
	Time     time.Time `json:"time"`
	TimeText string    `json:"time_text"`
	Location *Location `json:"location"`
}

// Location contains the location information for a food truck
type Location struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
	Type string `json:"type"`
}

// A Client communicates with stockholmfoodtrucks.nu
type Client struct {
	// URL is the url for stockholmfoodtrucks.nu
	URL *url.URL

	// User agent used for HTTP requests to stockholmfoodtrucks.nu
	UserAgent string

	// HTTP client used to communicate with stockholmfoodtrucks.nu
	httpClient *http.Client
}

// NewClient returns a new stockholmfoodtrucks client.
func NewClient(httpClients ...*http.Client) *Client {
	cloned := *http.DefaultClient
	httpClient := &cloned

	if len(httpClients) > 0 {
		httpClient = httpClients[0]
	}

	return &Client{
		URL: &url.URL{
			Scheme: Env("STOCKHOLM_FOOD_TRUCKS_URL_SCHEME", "http"),
			Host:   Env("STOCKHOLM_FOOD_TRUCKS_URL_HOST", "stockholmfoodtrucks.nu"),
		},
		UserAgent:  Env("STOCKHOLM_FOOD_TRUCKS_USER_AGENT", "stockholmfoodtrucks.go"),
		httpClient: httpClient,
	}
}

// NewRequest creates a new request to stockholmfoodtrucks.nu
func (c *Client) NewRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends a request and returns the response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Make sure to close the connection after replying to this request
	req.Close = true

	return c.httpClient.Do(req)
}

// NewDocument returns a goquery document based on stockholmfoodtrucks.nu
func (c *Client) NewDocument() (*goquery.Document, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return goquery.NewDocumentFromResponse(res)
}

// Get returns a slice of FoodTruck
func (c *Client) Get() ([]FoodTruck, error) {
	doc, err := c.NewDocument()
	if err != nil {
		return nil, err
	}

	return c.FoodTrucks(doc)
}

// FoodTrucks extracts slice of food trucks from goquery document
func (c *Client) FoodTrucks(doc *goquery.Document) ([]FoodTruck, error) {
	foodTrucks := []FoodTruck{}

	doc.Find(".trucks-list .truck").Each(func(i int, s *goquery.Selection) {
		truckName := s.Find(".truck-name").Text()

		post := s.Find(".posts .post").First()
		truckText := post.Find(".content").Text()
		truckTime, _ := time.Parse("2006-01-02 15:04", post.Find(".meta a").First().AttrOr("title", ""))
		truckTimeText := post.Find(".meta a").First().Text()

		foodTruck := FoodTruck{
			Name:     truckName,
			Text:     truckText,
			Time:     truckTime,
			TimeText: truckTimeText,
		}

		location := post.Find(".content .location").First()

		if id, exists := location.Attr("data-id"); exists {
			if n, exists := location.Attr("data-name"); exists {
				if t, exists := location.Attr("data-type"); exists {
					foodTruck.Location = &Location{
						ID:   id,
						Name: n,
						Type: t,
						Text: location.Text(),
					}
				}
			}
		}

		foodTrucks = append(foodTrucks, foodTruck)
	})

	return foodTrucks, nil
}

// Env returns a string from the ENV, or fallback variable
func Env(key, fallback string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}

	return fallback
}
