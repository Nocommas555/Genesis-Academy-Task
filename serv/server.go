package main

import (
	"net/http"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"
	"strings"
	"regexp"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"io"
)

type subscriber struct {
	Id 		int
	Email	string
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()    
	var res string
	// Ping test
	r.GET("/ping", func(c *gin.Context) {// 
		resp, err := http.Get("http://v6.exchangerate-api.com/v6/af0420066ad5474ea8b9fa25/pair/USD/UAH")
		
		if err != nil {
			fmt.Println("unable to query users: %v", err)
			return
		}

		if resp.StatusCode == http.StatusOK {
    		bodyBytes, err := io.ReadAll(resp.Body)

			if err != nil {
				return
			}

    		defer resp.Body.Close()
			var v map[string]interface{}
			json.Unmarshal(bodyBytes, &v)
			rawData, ok := v["conversion_rate"]
			if (!ok){
				res = "Technical troubles"
			}else{
				dateValue, _ := rawData.(float64)
				res = fmt.Sprintf("%f", dateValue) 
			}
		}else{
			res = "Technical troubles"
		}

		
		c.String(http.StatusOK, res)
	})

	
	r.POST("/api/subscribe", func(c *gin.Context) {

		bodyAsByteArray, _ := ioutil.ReadAll(c.Request.Body)
		email := string(bodyAsByteArray)
		// no sql injections pls
		email = strings.Replace(email, "`", "", -1)
		email = strings.Replace(email, "%40", "@", -1)
		email = strings.Replace(email, "email=", "", 1)
		

		//check if we got an actual email 
		matched, _ := regexp.MatchString(`^[^\s)"']+@[^\s)"']+\.[^\s)"']+$`, email)
		
		if(matched){
			rows, err := db.Query(context.Background(), `SELECT * FROM mailing_list_subcribers WHERE email='`+email+`'`)			
			
			if err != nil {
				fmt.Println("unable to query users: %v", err)
				return
			}
			
			arr, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[subscriber])
			
			defer rows.Close()
			
			fmt.Println(arr)
			if err != nil && arr.Email=="" {
				fmt.Println("unable to query users: %v", err)
				db.QueryRow(context.Background(), `INSERT INTO mailing_list_subcribers(email) VALUES('`+email+`')`)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Ok"})
			}else{
				c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Conflict"})
			}
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Malformed email"})
	})
	return r
}
var db *pgxpool.Pool
func main() {
	r := setupRouter()
	
	var err error
	db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("Successfully connected!")
		
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

