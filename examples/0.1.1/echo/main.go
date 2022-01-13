package main

import (
	"encoding/json"
	"fmt"
	"github.com/gojek/heimdall/v7/httpclient"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"time"

	auth "github.com/abbot/go-http-auth"
)

func secret(passwd string) func(user, realm string) string {
	return func(user, realm string) string {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
		return string(hashedPassword)
	}
}
func handle(url string) func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	return func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		Client := httpclient.NewClient(httpclient.WithHTTPTimeout(time.Second * 30))
		res, _ := Client.Get(fmt.Sprintf("http://%s/", url), nil)
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		_, _ = fmt.Fprintf(w, "Username:%s, Timebox: %s", r.Username, string(bodyBytes))
	}
}

func main() {
	url, _ := ioutil.ReadFile("/etc/configs/timebox")
	configBytes, _ := ioutil.ReadFile("/etc/configs/userconfigs")
	var conf struct {
		Username string
		Password string
	}
	_ = json.Unmarshal(configBytes, &conf)

	authenticator := auth.NewBasicAuthenticator("example.com", secret(conf.Password))
	http.HandleFunc("/", authenticator.Wrap(handle(string(url))))
	_ = http.ListenAndServe(":3000", nil)
}