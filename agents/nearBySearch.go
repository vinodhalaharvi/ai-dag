package agents

import (
	"ai-dag/config"
	"ai-dag/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type NearBySearchResponse struct {
	HtmlAttributions []interface{} `json:"html_attributions"`
	NextPageToken    string        `json:"next_page_token"`
	Results          []struct {
		BusinessStatus string `json:"business_status"`
		Geometry       struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			Viewport struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		Icon                string `json:"icon"`
		IconBackgroundColor string `json:"icon_background_color"`
		IconMaskBaseUri     string `json:"icon_mask_base_uri"`
		Name                string `json:"name"`
		OpeningHours        struct {
			OpenNow bool `json:"open_now"`
		} `json:"opening_hours,omitempty"`
		Photos []struct {
			Height           int      `json:"height"`
			HtmlAttributions []string `json:"html_attributions"`
			PhotoReference   string   `json:"photo_reference"`
			Width            int      `json:"width"`
		} `json:"photos"`
		PlaceId  string `json:"place_id"`
		PlusCode struct {
			CompoundCode string `json:"compound_code"`
			GlobalCode   string `json:"global_code"`
		} `json:"plus_code"`
		PriceLevel       int      `json:"price_level,omitempty"`
		Rating           float64  `json:"rating"`
		Reference        string   `json:"reference"`
		Scope            string   `json:"scope"`
		Types            []string `json:"types"`
		UserRatingsTotal int      `json:"user_ratings_total"`
		Vicinity         string   `json:"vicinity"`
	} `json:"results"`
	Status string `json:"status"`
}

type NearBySearchRequest struct {
	Location config.Location `json:"location"`
	Radius   int             `json:"radius"`
	Type     string          `json:"type"`
}

func NewNearBySearchRequest(
	location config.Location,
	radius int,
	Type string,
	key string,
) *NearBySearchRequest {
	return &NearBySearchRequest{
		Location: location,
		Radius:   radius,
		Type:     Type,
	}
}

type NearBySearch struct {
	*NearBySearchRequest
	*NearBySearchResponse
}

func NewNearBySearch(nearBySearchRequest *NearBySearchRequest) *NearBySearch {
	return &NearBySearch{NearBySearchRequest: nearBySearchRequest}
}

func (n *NearBySearch) toUrl() string {
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	if googleAPIKey == "" {
		log.Println("GOOGLE_API_KEY not set")
		return ""
	}
	return "https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=" +
		fmt.Sprintf("%f,%f", n.Location.Lat, n.Location.Lng) + "&radius=" +
		fmt.Sprintf("%d", n.Radius) + "&type=" + n.Type + "&key=" + googleAPIKey
}

func get(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}(resp.Body)

	return json.NewDecoder(resp.Body).Decode(target)
}

// Do Add Do method
func (n *NearBySearch) Do(
	config *config.DagConfig,
	agentId string,
	resultCh map[string]chan string,
	childResults map[string]string,
) {
	var response NearBySearchResponse
	err := get(n.toUrl(), &response)
	if err != nil {
		return
	}
	jsonStr := utils.ToPrettyJsonFromObject(response)
	resultCh[agentId] <- jsonStr
	close(resultCh[agentId])
}
