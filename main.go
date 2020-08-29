package main

import (
	// "net/http"
	"github.com/panatcmu/gofinal/crud"
	"github.com/panatcmu/gofinal/middleware"
	"github.com/gin-gonic/gin"
	// "log"
	"fmt"
	_ "github.com/lib/pq"
	
)


func main() {
	fmt.Println("customer service")
	r := gin.Default()
	r.Use(middleware.AuthMiddleware)
	r.POST("/customers", crud.CreateCustomers)
	r.GET("/customers/:id", crud.GetCustomer)
	r.GET("/customers", crud.GetCustomers)
	r.PUT("/customers/:id", crud.UpdateCustomers)
	r.DELETE("/customers/:id", crud.DeleteCustomer)
	r.Run(":2009")
	//run port ":2009"
}
