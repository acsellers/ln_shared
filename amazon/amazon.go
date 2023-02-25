package amazon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	RFAPIKey = "REPLACE_ME"

	PaperbackTypes = []string{"Paperback", "Mass Market Paperback", "Perfect Paperback", "Pocket Book"}
	HardcoverTypes = []string{"Hardcover", "Leather Bound", "Library Binding", "Flexibound"}
	DigitalTypes   = []string{"Kindle", "Kindle & Comixology", "Digital"}
	AudiobookTypes = []string{"Audiobook", "Audible Audiobook", "Audio CD", "MP3 CD"}

	cacheList = map[string]string{}
	cacheMtx  sync.RWMutex
)

type ProductData struct {
	RequestInfo struct {
		Success                bool      `json:"success"`
		CreditsUsed            int       `json:"credits_used"`
		CreditsUsedThisRequest int       `json:"credits_used_this_request"`
		CreditsRemaining       int       `json:"credits_remaining"`
		CreditsResetAt         time.Time `json:"credits_reset_at"`
	} `json:"request_info"`
	RequestParameters struct {
		AmazonDomain string `json:"amazon_domain"`
		Type         string `json:"type"`
		Gtin         string `json:"gtin"`
	} `json:"request_parameters"`
	RequestMetadata struct {
		CreatedAt      time.Time `json:"created_at"`
		ProcessedAt    time.Time `json:"processed_at"`
		TotalTimeTaken float64   `json:"total_time_taken"`
		AmazonURL      string    `json:"amazon_url"`
	} `json:"request_metadata"`
	Product struct {
		Title       string `json:"title"`
		SearchAlias struct {
			Title string `json:"title"`
			Value string `json:"value"`
		} `json:"search_alias"`
		Keywords         string           `json:"keywords"`
		KeywordsList     []string         `json:"keywords_list"`
		Asin             string           `json:"asin"`
		Link             string           `json:"link"`
		SellOnAmazon     bool             `json:"sell_on_amazon"`
		Variants         []ProductVariant `json:"variants"`
		VariantAsinsFlat string           `json:"variant_asins_flat"`
		Authors          []struct {
			Name string `json:"name"`
			Link string `json:"link"`
			Asin string `json:"asin"`
		} `json:"authors"`
		Format     string `json:"format"`
		Categories []struct {
			Name       string `json:"name"`
			Link       string `json:"link"`
			CategoryID string `json:"category_id"`
		} `json:"categories"`
		CategoriesFlat string `json:"categories_flat"`
		SubTitle       struct {
			Text string `json:"text"`
			Link string `json:"link"`
		} `json:"sub_title"`
		MarketplaceID   string  `json:"marketplace_id"`
		Rating          float64 `json:"rating"`
		RatingBreakdown struct {
			FiveStar struct {
				Percentage float32 `json:"percentage"`
				Count      int     `json:"count"`
			} `json:"five_star"`
			FourStar struct {
				Percentage float32 `json:"percentage"`
				Count      int     `json:"count"`
			} `json:"four_star"`
			ThreeStar struct {
				Percentage float32 `json:"percentage"`
				Count      int     `json:"count"`
			} `json:"three_star"`
			TwoStar struct {
				Percentage float32 `json:"percentage"`
				Count      int     `json:"count"`
			} `json:"two_star"`
			OneStar struct {
				Percentage float32 `json:"percentage"`
				Count      int     `json:"count"`
			} `json:"one_star"`
		} `json:"rating_breakdown"`
		RatingsTotal     int    `json:"ratings_total"`
		BookDescription  string `json:"book_description"`
		EditorialReviews []struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		} `json:"editorial_reviews"`
		EditorialReviewsFlat string `json:"editorial_reviews_flat"`
		MainImage            struct {
			Link string `json:"link"`
		} `json:"main_image"`
		Images []struct {
			Link string `json:"link"`
		} `json:"images"`
		ImagesCount int    `json:"images_count"`
		ImagesFlat  string `json:"images_flat"`
		IsBundle    bool   `json:"is_bundle"`
		Attributes  []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"attributes"`
		TopReviews []struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Body     string `json:"body"`
			BodyHTML string `json:"body_html"`
			Link     string `json:"link,omitempty"`
			Rating   int    `json:"rating"`
			Date     struct {
				Raw string    `json:"raw"`
				Utc time.Time `json:"utc"`
			} `json:"date"`
			Profile struct {
				Name string `json:"name"`
				Link string `json:"link"`
				ID   string `json:"id"`
			} `json:"profile,omitempty"`
			VineProgram      bool   `json:"vine_program"`
			VerifiedPurchase bool   `json:"verified_purchase"`
			ReviewCountry    string `json:"review_country"`
			IsGlobalReview   bool   `json:"is_global_review"`
			HelpfulVotes     int    `json:"helpful_votes,omitempty"`
			Profile0         struct {
				Name string `json:"name"`
			} `json:"profile,omitempty"`
			Profile1 struct {
				Name string `json:"name"`
			} `json:"profile,omitempty"`
			Profile2 struct {
				Name string `json:"name"`
			} `json:"profile,omitempty"`
			Profile3 struct {
				Name  string `json:"name"`
				Image string `json:"image"`
			} `json:"profile,omitempty"`
			Profile4 struct {
				Name string `json:"name"`
			} `json:"profile,omitempty"`
		} `json:"top_reviews"`
		BuyboxWinner struct {
			MaximumOrderQuantity struct {
				Value       int  `json:"value"`
				HardMaximum bool `json:"hard_maximum"`
			} `json:"maximum_order_quantity"`
			SecondaryBuybox struct {
				OfferID string `json:"offer_id"`
				Caption string `json:"caption"`
				Price   struct {
					Symbol   string  `json:"symbol"`
					Value    float64 `json:"value"`
					Currency string  `json:"currency"`
					Raw      string  `json:"raw"`
				} `json:"price"`
				Availability struct {
					Raw string `json:"raw"`
				} `json:"availability"`
			} `json:"secondary_buybox"`
			OfferID        string `json:"offer_id"`
			NewOffersCount int    `json:"new_offers_count"`
			NewOffersFrom  struct {
				Value    float64 `json:"value"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
				Raw      string  `json:"raw"`
			} `json:"new_offers_from"`
			UsedOffersCount int `json:"used_offers_count"`
			UsedOffersFrom  struct {
				Value    float64 `json:"value"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
				Raw      string  `json:"raw"`
			} `json:"used_offers_from"`
			IsPrime       bool `json:"is_prime"`
			IsAmazonFresh bool `json:"is_amazon_fresh"`
			Condition     struct {
				IsNew bool `json:"is_new"`
			} `json:"condition"`
			Availability struct {
				Type         string `json:"type"`
				Raw          string `json:"raw"`
				DispatchDays int    `json:"dispatch_days"`
				StockLevel   int    `json:"stock_level"`
			} `json:"availability"`
			Fulfillment struct {
				Type             string `json:"type"`
				StandardDelivery struct {
					Date string `json:"date"`
					Name string `json:"name"`
				} `json:"standard_delivery"`
				FastestDelivery struct {
					Date string `json:"date"`
					Name string `json:"name"`
				} `json:"fastest_delivery"`
				IsSoldByAmazon          bool `json:"is_sold_by_amazon"`
				IsFulfilledByAmazon     bool `json:"is_fulfilled_by_amazon"`
				IsFulfilledByThirdParty bool `json:"is_fulfilled_by_third_party"`
				IsSoldByThirdParty      bool `json:"is_sold_by_third_party"`
			} `json:"fulfillment"`
			Price struct {
				Symbol   string  `json:"symbol"`
				Value    float64 `json:"value"`
				Currency string  `json:"currency"`
				Raw      string  `json:"raw"`
			} `json:"price"`
			Shipping struct {
				Raw string `json:"raw"`
			} `json:"shipping"`
		} `json:"buybox_winner"`
		MoreBuyingChoices []struct {
			Price struct {
				Symbol   string  `json:"symbol"`
				Value    float64 `json:"value"`
				Currency string  `json:"currency"`
				Raw      string  `json:"raw"`
			} `json:"price"`
			SellerName   string `json:"seller_name"`
			SellerLink   string `json:"seller_link"`
			FreeShipping bool   `json:"free_shipping,omitempty"`
			Position     int    `json:"position"`
		} `json:"more_buying_choices"`
		Specifications []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"specifications"`
		SpecificationsFlat string `json:"specifications_flat"`
		BestsellersRank    []struct {
			Category string `json:"category"`
			Rank     int    `json:"rank"`
			Link     string `json:"link"`
		} `json:"bestsellers_rank"`
		PublicationDate     string `json:"publication_date"`
		Publisher           string `json:"publisher"`
		Isbn10              string `json:"isbn_10"`
		Isbn13              string `json:"isbn_13"`
		Language            string `json:"language"`
		Weight              string `json:"weight"`
		BestsellersRankFlat string `json:"bestsellers_rank_flat"`
	} `json:"product"`
	FrequentlyBoughtTogether struct {
		TotalPrice struct {
			Symbol   string  `json:"symbol"`
			Value    float64 `json:"value"`
			Currency string  `json:"currency"`
			Raw      string  `json:"raw"`
		} `json:"total_price"`
		Products []struct {
			Asin  string `json:"asin"`
			Title string `json:"title"`
			Link  string `json:"link"`
			Price struct {
				Symbol   string  `json:"symbol"`
				Value    float64 `json:"value"`
				Currency string  `json:"currency"`
				Raw      string  `json:"raw"`
			} `json:"price"`
			Image string `json:"image,omitempty"`
		} `json:"products"`
	} `json:"frequently_bought_together"`
	AlsoBought []struct {
		Title        string  `json:"title"`
		Asin         string  `json:"asin"`
		Link         string  `json:"link"`
		Image        string  `json:"image"`
		Rating       float32 `json:"rating"`
		RatingsTotal float32 `json:"ratings_total"`
		Price        struct {
			Symbol   string  `json:"symbol"`
			Value    float64 `json:"value"`
			Currency string  `json:"currency"`
			Raw      string  `json:"raw"`
		} `json:"price"`
	} `json:"also_bought"`
}
type ProductVariant struct {
	Asin             string `json:"asin"`
	Link             string `json:"link"`
	IsCurrentProduct bool   `json:"is_current_product"`
	Title            string `json:"title"`
	Price            struct {
		Symbol   string  `json:"symbol"`
		Value    float64 `json:"value"`
		Currency string  `json:"currency"`
		Raw      string  `json:"raw"`
	} `json:"price"`
}

func (pd ProductData) LookupVariant(titles ...string) (ProductVariant, bool) {
	for _, title := range titles {
		for _, v := range pd.Product.Variants {
			if v.Title == title {
				return v, true
			}
		}
	}
	return ProductVariant{}, false
}
func init() {
	cacheList = make(map[string]string)
	folder := fmt.Sprintf("amazon/%s/", time.Now().Format("2006-01"))
	matches, _ := filepath.Glob(folder + "/*.json")
	for _, match := range matches {
		id := strings.Replace(match, folder, "", 1)
		id = strings.Replace(id, ".json", "", 1)
		cacheList[id] = match
	}
	f, _ := os.Open("amazon/missing.json")
	missing := []string{}
	json.NewDecoder(f).Decode(&missing)
	f.Close()
	for _, id := range missing {
		cacheList[id] = "missing"
	}
}
func SaveMissing(id string) {
	f, _ := os.Create(fmt.Sprintf("amazon/%s/missing.json", time.Now().Format("2006-01")))
	missing := []string{}
	for k, v := range cacheList {
		if v == "missing" {
			missing = append(missing, k)
		}
	}
	json.NewEncoder(f).Encode(missing)
	f.Close()
}

func CacheData(id string, pd ProductData) {
	filename := fmt.Sprintf("amazon/%s/%s.json", time.Now().Format("2006-01"), id)
	cacheMtx.Lock()
	cacheList[id] = filename
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Create: ", err)
	}

	err = json.NewEncoder(f).Encode(pd)
	if err != nil {
		log.Fatal("Encode: ", err)
	}
	f.Close()
	cacheMtx.Unlock()
}

func RetrieveASIN(id string) ProductData {
	cacheMtx.RLock()
	if cached, ok := cacheList[id]; ok {
		cacheMtx.RUnlock()
		fmt.Println("Cached: ", id)
		if cached == "missing" {
			return ProductData{}
		}
		f, _ := os.Open(cached)
		pd := ProductData{}
		err := json.NewDecoder(f).Decode(&pd)
		if err != nil {
			fmt.Println("ID: ", id)
			log.Fatal("Decode: ", err)
		}
		return pd
	}
	cacheMtx.RUnlock()

	u := fmt.Sprintf(
		"https://api.rainforestapi.com/request?api_key=%s&amazon_domain=amazon.com&asin=%s&type=product",
		RFAPIKey,
		id,
	)
	fmt.Println("Retrieving: ", id)
	pd, err := Get(u)
	if err != nil {
		fmt.Println("ID: ", id)
		log.Fatal("Get: ", err)
	}
	if pd.Product.Asin == "" {
		fmt.Println("Not Found: ", id)
		return ProductData{}
	}
	CacheData(id, pd)
	return pd
}
func RetrieveGTIN(id string) ProductData {
	cacheMtx.RLock()
	if cached, ok := cacheList[id]; ok {
		cacheMtx.RUnlock()
		fmt.Println("Cached: ", id)
		if cached == "missing" {
			return ProductData{}
		}
		f, _ := os.Open(cached)
		pd := ProductData{}
		err := json.NewDecoder(f).Decode(&pd)
		if err != nil {
			fmt.Println("ID: ", id)
			log.Fatal("Decode: ", err)
		}
		return pd
	}
	cacheMtx.RUnlock()

	u := fmt.Sprintf(
		"https://api.rainforestapi.com/request?api_key=%s&amazon_domain=amazon.com&type=product&gtin=%s",
		RFAPIKey,
		id,
	)
	fmt.Println("Retrieving: ", id)
	pd, err := Get(u)
	if err != nil {
		fmt.Println("ID: ", id)
		log.Fatal("Get: ", err)
	}
	if pd.Product.Asin == "" {
		fmt.Println("Not Found: ", id)
		cacheList[id] = "missing"
		return ProductData{}
	}
	CacheData(id, pd)
	return pd
}

func Get(url string) (ProductData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return ProductData{}, err
	}
	defer resp.Body.Close()
	var pd ProductData
	err = json.NewDecoder(resp.Body).Decode(&pd)
	return pd, err
}
