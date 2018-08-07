package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "psoberanes:Paul1990+@tcp(127.0.0.1:3306)/gotest")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	type Person struct {
		ID        int
		FirstName string
		LastName  string
	}
	router := gin.Default()

	// API handlers

	// GET a person detail
	router.GET("/person/:id", func(c *gin.Context) {
		var (
			person Person
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("SELECT id, first_name, last_name FROM person WHERE id = ?;", id)
		err = row.Scan(&person.ID, &person.FirstName, &person.LastName)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": person,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET all persons
	router.GET("/persons", func(c *gin.Context) {
		var (
			person  Person
			persons []Person
		)
		rows, err := db.Query("SELECT * FROM person;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&person.ID, &person.FirstName, &person.LastName)
			persons = append(persons, person)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": persons,
			"count":  len(persons),
		})
	})

	// POST new person details
	router.POST("/person", func(c *gin.Context) {
		var buffer bytes.Buffer
		firstName := c.PostForm("first_name")
		lastName := c.PostForm("last_name")
		stmt, err := db.Prepare("INSERT INTO person (first_name, last_name) VALUES (?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(firstName, lastName)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(firstName)
		buffer.WriteString(" ")
		buffer.WriteString(lastName)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// DELETE resources
	router.DELETE("/person", func(c *gin.Context) {
		id := c.Query("id")
		stmt, err := db.Prepare("DELETE FROM person WHERE id = ?;")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})

	})

	// PUT update a person details
	router.PUT("/person", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.Query("id")
		firstName := c.PostForm("first_name")
		lastName := c.PostForm("last_name")
		stmt, err := db.Prepare("UPDATE person SET first_name = ?, last_name = ? WHERE id = ?;")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(firstName, lastName, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(firstName)
		buffer.WriteString(" ")
		buffer.WriteString(lastName)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	// End API handlers

	// Server runs on 3000 port
	router.Run(":3000")
}
