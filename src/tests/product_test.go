package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestProduct_CRUDWithRelations(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "prod_crud@test.com", "password123")

	// Create a tag
	w := performRequest("POST", "/v1/tag/",
		map[string]string{"name": "Product Tag"}, accessToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("create tag: expected 201, got %d", w.Code)
	}

	// Get tag ID
	w = performRequest("GET", "/v1/tag/", nil, "")
	tagResp := parseJSON(w)
	tags := tagResp["response"].([]interface{})
	tagItem := tags[len(tags)-1].(map[string]interface{})
	tagID := tagItem["_id"].(string)

	// Create a category
	w = performRequest("POST", "/v1/category/",
		map[string]string{"name": "Product Category"}, accessToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("create category: expected 201, got %d", w.Code)
	}

	// Get category ID
	w = performRequest("GET", "/v1/category/", nil, "")
	catResp := parseJSON(w)
	cats := catResp["response"].([]interface{})
	catItem := cats[len(cats)-1].(map[string]interface{})
	catID := catItem["_id"].(string)

	// CREATE PRODUCT with relations
	productBody := map[string]interface{}{
		"name":       "Test Product",
		"brand":      "Test Brand",
		"categories": []string{catID},
		"tags":       []string{tagID},
		"summary":    "A test product",
	}
	w = performRequest("POST", "/v1/product/", productBody, accessToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("create product: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// LIST (uses $lookup aggregation)
	w = performRequest("GET", "/v1/product/", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("list products: expected 200, got %d", w.Code)
	}

	prodResp := parseJSON(w)
	products := prodResp["response"].([]interface{})
	if len(products) == 0 {
		t.Fatal("expected at least one product")
	}

	prodItem := products[len(products)-1].(map[string]interface{})
	prodID := prodItem["_id"].(string)

	// Verify product has populated tags (from $lookup)
	if prodTags, ok := prodItem["tags"].([]interface{}); ok && len(prodTags) > 0 {
		firstTag, ok := prodTags[0].(map[string]interface{})
		if ok && firstTag["name"] != "Product Tag" {
			t.Errorf("expected populated tag name 'Product Tag', got '%v'", firstTag["name"])
		}
	}

	// FIND (uses $lookup aggregation)
	w = performRequest("GET", fmt.Sprintf("/v1/product/%s", prodID), nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("find product: expected 200, got %d", w.Code)
	}

	findResp := parseJSON(w)
	product := findResp["response"].(map[string]interface{})
	if product["name"] != "Test Product" {
		t.Errorf("expected 'Test Product', got '%v'", product["name"])
	}
	if product["brand"] != "Test Brand" {
		t.Errorf("expected 'Test Brand', got '%v'", product["brand"])
	}

	// UPDATE
	w = performRequest("PATCH", fmt.Sprintf("/v1/product/%s", prodID),
		map[string]interface{}{
			"name":       "Updated Product",
			"brand":      "Updated Brand",
			"categories": []string{catID},
		}, accessToken)
	if w.Code != http.StatusOK {
		t.Fatalf("update product: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// DELETE
	w = performRequest("DELETE", fmt.Sprintf("/v1/product/%s", prodID), nil, accessToken)
	if w.Code != http.StatusOK {
		t.Fatalf("delete product: expected 200, got %d", w.Code)
	}

	// Verify delete
	w = performRequest("GET", fmt.Sprintf("/v1/product/%s", prodID), nil, "")
	if w.Code != http.StatusNotFound {
		t.Errorf("find after delete: expected 404, got %d", w.Code)
	}
}

func TestProduct_FindInvalidID(t *testing.T) {
	w := performRequest("GET", "/v1/product/invalid-id", nil, "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestProduct_CreateWithoutAuth(t *testing.T) {
	w := performRequest("POST", "/v1/product/",
		map[string]interface{}{"name": "No Auth", "brand": "Test", "categories": []string{}}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestProduct_ListEmpty(t *testing.T) {
	w := performRequest("GET", "/v1/product/", nil, "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	if _, ok := resp["response"].([]interface{}); !ok {
		t.Fatalf("expected array response, got %T", resp["response"])
	}
}
