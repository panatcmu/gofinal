package crud

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"fmt"
)

type Customer struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Status string `json:"status"`
}

type Response struct {
	Message string `json:"message"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB connected")

	
}

func CreateCustomers(c *gin.Context) {
	customer := Customer{}
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	createTb := `
    CREATE TABLE IF NOT EXISTS customers (
        id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
        status TEXT
    );
    `
	_, err := db.Exec(createTb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println("create table success or existed")
	row := db.QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id", customer.Name, customer.Email, customer.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	customer.ID = id
	c.JSON(http.StatusCreated, customer)
}

func GetCustomer(c *gin.Context) {
	
	paramId := c.Param("id")
	fmt.Println(paramId)
	stmt, err := db.Prepare("Select id , name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	rowId := paramId
	row := stmt.QueryRow(rowId)
	var id int
	var name, email, status string

	err = row.Scan(&id, &name, &email, &status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	customer := Customer{id, name, email, status}
	c.JSON(http.StatusOK, customer)
}

func GetCustomers(c *gin.Context) {
	// customers := []Customer{}
	var customers []Customer
	stmt, err := db.Prepare("Select id , name, email, status FROM customers ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for rows.Next() {
		var id int
		var name, email, status string
		err = rows.Scan(&id, &name, &email, &status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		customer := Customer{id, name, email, status}
		customers = append (customers, customer)
	}
	c.JSON(http.StatusOK, customers)
}

func UpdateCustomers(c *gin.Context) {
	paramId := c.Param("id")
	fmt.Println(paramId)
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	row := stmt.QueryRow(paramId)
	customer := &Customer{}
	err = row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := c.ShouldBindJSON(customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err = db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if _, err := stmt.Exec(paramId, customer.Name, customer.Email, customer.Status); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, customer)
}

func DeleteCustomer(c *gin.Context) {
	paramId := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _,err := stmt.Exec(paramId); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	t := Response{"customer deleted"}
	c.JSON(http.StatusOK, t)
}
