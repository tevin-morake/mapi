package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Rapidtrade/rapi/tools"
	"github.com/tevin-morake/mapi/controllers"
)

var (
	port        string
	version     = "1"
	showversion bool
	stdout      bool
	t           *tools.Tools
)

func init() {
	flag.StringVar(&port, "port", "80", "Use port 80 in production")
	flag.BoolVar(&showversion, "version", false, "Show current version")
	flag.Parse()

	if showversion {
		fmt.Printf("Current Mail API Version : %s \n", version)
		os.Exit(0)
	}

	//Initialize a tools service we can use to connect to the db throughout the app
	t, err = tools.NewTools(stdout, "logs", "logs", "mapi")
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/email", handleCorsMiddleWare(PostEmail))
	http.Handle("/", mux)

	fmt.Printf("Listening on port %s\n", port)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  900 * time.Second,
		WriteTimeout: 900 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

//PostEmail runs a controller method that sends an email to an address specified in the request body
func PostEmail(resp http.ResponseWriter, req *http.Request) {
	controllers.SendEmail(resp, req, t)
}

//handleCorsMiddleWare is used to handle preflight CORS requests
func handleCorsMiddleWare(h http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get("Origin"); origin != "" {
			res.Header().Set("Access-Control-Allow-Origin", origin)
			res.Header().Set("Access-Control-Allow-Credentials", "true")
			res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			res.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

		if strings.ToLower(req.Method) == "options" {
			fmt.Println("CORS Preflight successful")
			res.WriteHeader(http.StatusOK)
			return
		}
		// Check if  request has basic auth . If not, deny request further access
		userid, password, ok := req.BasicAuth()
		if !ok || userid == "" || password == "" {
			http.Error(res, "401 - unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(res, req)
	}
}
