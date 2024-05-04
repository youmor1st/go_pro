package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"shop/models"
	"strconv"
)

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := models.GetOrdersByUserID(currentUser.ID)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
func GetIDOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := mux.Vars(r)["order_id"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := models.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Failed to fetch order", http.StatusInternalServerError)
		return
	}
	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var orderRequest struct {
		ProductIDs []int `json:"product_ids"`
	}
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var products []*models.Product
	for _, productID := range orderRequest.ProductIDs {
		product, err := models.GetProductByID(productID)
		if err != nil {
			http.Error(w, "Failed to get product information", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	var totalAmount float64
	for _, product := range products {
		totalAmount += product.Price
	}

	order := &models.Order{
		UserID:      currentUser.ID,
		TotalAmount: totalAmount,
		Status:      "created",
	}
	err = models.CreateOrder(order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	orderID := order.ID

	for _, product := range products {
		orderItem := &models.OrderItem{
			OrderID:   orderID,
			ProductID: product.ID,
			Quantity:  1,
			Price:     product.Price,
		}
		err = models.CreateOrderItem(orderItem)
		if err != nil {
			http.Error(w, "Failed to add product to order", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := mux.Vars(r)["order_id"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var updateRequest struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := models.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if currentUser.ID != order.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	order.Status = updateRequest.Status
	err = models.UpdateOrder(order)
	if err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := mux.Vars(r)["order_id"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := models.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if currentUser.ID != order.UserID {
		http.Error(w, "Forbidden: You are not allowed to delete this order", http.StatusForbidden)
		return
	}

	err = models.DeleteOrder(orderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete order: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
