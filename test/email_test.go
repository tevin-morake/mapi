package main

import "testing"
var (
	authuser = "postgres"
	authpassword = "postgres"
)

//TODO: Expand more on this test
TestSendEmail(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	userpwd := authuser + ":" + authpassword
	auth = "Basic " + b64.StdEncoding.EncodeToString([]byte(userpwd))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	tl, _ := tools.NewTools(false, "logs", "mapi")

	testmail := `{
		"email": "rory@sharklasers.com",
		"subject": "certificate authority",
		"body": "We detected a new certificate created in your name. Is this you ? If not , let's stop the fraud "}`
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(testmail)))
	controllers.SendEmail(w, req, tl)
	if w.Code != http.StatusOK {
		t.Errorf("Expected req to bring back response %v, instead got %s", http.StatusOK, w.Code)
		return false
	}

	return true
}