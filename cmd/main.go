package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"shop/handlers"
	"shop/models"
)

func main() {
	models.ConnectDB()
	defer models.CloseDB()

	router := mux.NewRouter()

	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/profile", handlers.GetProfileHandler).Methods("GET")
	router.HandleFunc("/profile/update", handlers.UpdateProfileHandler).Methods("PUT")
	router.HandleFunc("/profile/delete", handlers.DeleteProfileHandler).Methods("DELETE")
	router.HandleFunc("/profile/{username}", handlers.GetUserProfileHandler).Methods("GET")
	router.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET")
	router.HandleFunc("/products/add", handlers.AddProduct).Methods("POST")
	router.HandleFunc("/products/{id}/update", handlers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}/delete", handlers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/myproducts", handlers.GetMyProducts).Methods("GET")
	router.HandleFunc("/cart", handlers.GetCartHandler).Methods("GET")
	router.HandleFunc("/cart/add/{product_id}", handlers.AddProductToCartHandler).Methods("POST")
	router.HandleFunc("/cart/update/{product_id}", handlers.UpdateCartItemHandler).Methods("PUT")
	router.HandleFunc("/cart/remove/{product_id}", handlers.RemoveProductFromCartHandler).Methods("DELETE")
	router.HandleFunc("/orders", handlers.GetOrdersHandler).Methods("GET")
	router.HandleFunc("/orders/{order_id}", handlers.GetIDOrderHandler).Methods("GET")
	router.HandleFunc("/orders/create", handlers.CreateOrderHandler).Methods("POST")
	router.HandleFunc("/orders/update/{order_id}", handlers.UpdateOrderHandler).Methods("PUT")
	router.HandleFunc("/orders/remove/{order_id}", handlers.DeleteOrderHandler).Methods("DELETE")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(`:8080`, router))
}
