package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"shop/models"
	"strconv"
)

func GetCartHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cart, err := models.GetCartByUserID(currentUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func AddProductToCartHandler(w http.ResponseWriter, r *http.Request) {
	productIDStr := mux.Vars(r)["product_id"]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	quantityStr := r.FormValue("quantity")
	quantity := 1
	if quantityStr != "" {
		quantity, err = strconv.Atoi(quantityStr)
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}
	}

	cartItem := &models.CartItem{
		UserID:    currentUser.ID,
		ProductID: productID,
		Quantity:  quantity,
	}

	err = models.AddProductToCart(cartItem)
	if err != nil {
		http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	productIDStr := mux.Vars(r)["product_id"]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updateRequest struct {
		Quantity int `json:"quantity"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updateRequest.Quantity < 1 {
		updateRequest.Quantity = 1
	}

	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	carts, err := models.GetCartByUserID(currentUser.ID)
	if err != nil {
		http.Error(w, "Failed to get user's cart", http.StatusInternalServerError)
		return
	}

	var cartItem *models.CartItem
	for _, cart := range carts {
		if item, err := models.GetCartItemByProductID(cart.ID, productID); err == nil {
			cartItem = item
			break
		}
	}

	if cartItem == nil {
		http.Error(w, "Product not found in cart", http.StatusNotFound)
		return
	}

	cartItem.Quantity = updateRequest.Quantity
	err = models.UpdateCartItem(cartItem)
	if err != nil {
		http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RemoveProductFromCartHandler(w http.ResponseWriter, r *http.Request) {
	productIDStr := mux.Vars(r)["product_id"]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cart, err := models.GetCartByUserID(currentUser.ID)
	if err != nil {
		http.Error(w, "Failed to get user's cart", http.StatusInternalServerError)
		return
	}

	var targetIndex int = -1

	for i, item := range cart {
		if item.ProductID == productID {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		http.Error(w, "Product not found in cart", http.StatusNotFound)
		return
	}

	updatedCart := append(cart[:targetIndex], cart[targetIndex+1:]...)

	err = models.UpdateCartByUserID(currentUser.ID, updatedCart)
	if err != nil {
		http.Error(w, "Failed to remove product from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
