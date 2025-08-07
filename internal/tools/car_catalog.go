package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/repository"
	"github.com/aryadhira/otogenius-agent/internal/storages"
)

func GetCarCatalog(brand, model, category, transmission string, production_year int, price float64) (string, error) {
	db, err := storages.NewDB()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	carInfo := repository.NewCarRepo(ctx, db)

	filter := make(map[string]any)
	if brand != "" {
		filter["brand"] = brand
	}
	if model != "" {
		filter["model"] = model
	}
	if category != "" {
		filter["category"] = category
	}
	if transmission != "" {
		filter["transmission"] = transmission
	}
	if production_year > 0 {
		filter["production_year"] = production_year
	}
	if price > 0 {
		filter["price"] = price
	}

	infos, err := carInfo.GetCarData(filter)
	if err != nil {
		return "", err
	}

	var carList strings.Builder
	for _, each := range infos {
		str := fmt.Sprintf("Brand: %s, Model: %s, Production Year: %v, Transmission: %s, Fuel: %s, Price: %v\n", each.Brand, each.Model, each.ProductionYear, each.Transmission, each.Fuel, int(each.Price))
		carList.WriteString(str)
	}

	return carList.String(), nil
}

func GetCarCatalogToolDescription() models.Function {
	return models.Function{
		Name:        "get_car_catalog",
		Description: "retrieve latest list of used car catalog",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"brand": map[string]interface{}{
					"type":        "string",
					"description": "The car brand eg Toyota,Honda etc. can be multiple by comma separated, can be empty pass with empty string",
				},
				"model": map[string]interface{}{
					"type":        "string",
					"description": "The car brand model eg Civic,Corolla etc. can be multiple by comma separated, can be empty pass with empty string",
				},
				"category": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"Sedan", "SUV", "MPV", "Hatchback"},
					"description": "The car category",
				},
				"transmission": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"", "Automatic", "Manual"},
					"description": "The car transmission type",
				},
				"production_year": map[string]interface{}{
					"type":        "integer",
					"description": "Car production year can be 0 as empty parameter",
				},
				"price": map[string]interface{}{
					"type":        "number",
					"description": "Car price must be > 20000000",
				},
			},
			"required": []string{"brand", "model", "category", "transmission", "production_year", "price"},
		},
	}
}

func ParseCarCatalogToolParameter(fnType reflect.Type, toolCall models.FunctionCall) ([]reflect.Value, error) {
	var input []reflect.Value
	fmt.Println("arguments:", toolCall.Function.Arguments)

	var arg map[string]interface{}
	err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal tool arguments: %w", err)
	}

	keyParam := []string{"brand", "model", "category", "transmission", "production_year", "price"}

	for i := 0; i < fnType.NumIn(); i++ {
		key := keyParam[i]
		if i < 4 {
			if arg[key] != nil {
				val := arg[key].(string)
				input = append(input, reflect.ValueOf(val))
			} else {
				emptyStr := ""
				input = append(input, reflect.ValueOf(emptyStr))
			}
		} else if i == 4 {
			if arg[key] != nil {
				val := arg[key].(float64)
				input = append(input, reflect.ValueOf(int(val)))
			}
		} else if i == 5 {
			if arg[key] != nil {
				val := arg[key].(float64)
				input = append(input, reflect.ValueOf(val))
			}
		}

	}

	return input, nil
}
