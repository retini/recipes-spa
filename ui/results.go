package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ApiResult struct {
	Hits []Recipe
}

type Recipe struct {
	Recipe struct {
		Label       string
		Image       string
		Source      string
		Ingredients []struct {
			Text string
		}
		TotalNutrients struct {
			EnercKcal struct {
				Quantity float64
			} `json:"ENERC_KCAL"`
		} `json:"totalNutrients"`
	}
}

func (p *Page) Results(res http.ResponseWriter, req *http.Request) Tags {
	ingredient := req.URL.Query().Get("ingredient")
	if ingredient == "" {
		log.Print("must pass a valid ingredient")
		return nil
	}
	q := req.URL.Query().Get("quantity")
	if q == "" {
		log.Print("must pass a valid quantity")
		return nil
	}
	quantity, err := strconv.Atoi(q)
	if err != nil {
		log.Print("invalid quantity: ", err)
		return nil
	}
	if quantity < 1 || quantity > 20 {
		log.Print("quantity must be between 1 and 20 (inclusive)")
		return nil
	}
	apiResult, err := FindRecipes(ingredient, quantity)
	if err != nil {
		log.Print("call to recipes Api produced following error: ", err)
		return nil
	}

	recipes := apiResult.Hits

	return map[string]interface{}{
		"recipes": recipes,
	}
}

func FindRecipes(i string, q int) (*ApiResult, error) {
	res, err := http.Get(fmt.Sprintf("https://api.edamam.com/search?q=%s&from=0&to=%d&app_id=0ae29fa7&app_key=2cd7ee8ee8691774c8f018c615dcf8f2", i, q))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	var results ApiResult
	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
