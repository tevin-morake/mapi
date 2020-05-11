package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/tevin-morake/mapi/tools"
)

//EmailBody holds  info about the email to be sent
type EmailBody struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

//SendEmail is responsible for sending out emails to email address supplied to it
func SendEmail(w http.ResponseWriter, r *http.Request, t *tools.Tools) {
	// Get body into a byte array
	body, ioerr := ioutil.ReadAll(io.LimitReader(r.Body, 104857600))
	if ioerr != nil {
		http.Error(w, fmt.Sprintf("Error reading body %s", ioerr.Error()), http.StatusBadRequest)
		return
	}

	//handle errors in closing request body
	if err := r.Body.Close(); err != nil {
		http.Error(w, fmt.Sprintf("Error closing body %s", err.Error()), http.StatusBadRequest)
		return
	}

	//transform body into map[string]string
	var emailbody EmailBody
	if err := json.Unmarshal(body, &emailbody); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling input into map : %s", err.Error()), http.StatusBadRequest)
		return
	}

	//Email a user
	from := mail.NewEmail("noreply user", "noreply@morakeholdings.com")
	to := mail.NewEmail("system user", emailbody.Email)

	message := mail.NewSingleEmail(from, emailbody.Subject, to, emailbody.Body, "")
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending email to %s : %s", emailbody.Email, err.Error()), http.StatusBadRequest)
		return
	} else {
		//connect to Postgres db
		if err := t.OpenDB(w); err != nil {
			http.Error(w, fmt.Sprintf("Error opening database : %s", err.Error()), http.StatusBadRequest)
			return
		}

		//execute sql
		sqlStatement := "INSERT INTO logs (logtype, email, subject, body) VALUES ($1, $2, $3, $4)"
		statementValues := []interface{"INFO",emailbody.Email, emailbody.Subject, emailbody.Body}
		
		if err = t.ExecuteSQL(sqlStatement, statementValues); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Send response to user
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}
