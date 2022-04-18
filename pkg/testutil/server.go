package testutil

import (
	"fmt"
	"net/http/httptest"
	"net/url"

	"github.com/Schnides123/wordle-vs/pkg/data"
	"github.com/Schnides123/wordle-vs/pkg/endpoints"
	"github.com/gorilla/mux"
)

var srv *httptest.Server
var serverURL *url.URL

func StartTestServer() bool {
	r := mux.NewRouter()
	endpoints.SetupRoutes(r)
	srv = httptest.NewServer(r)
	var err error
	serverURL, err = url.Parse(srv.URL)

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func StopTestServer(_ bool) {
	srv.CloseClientConnections()
	srv.Close()

	// reset in-memory sessions
	endpoints.ResetBlitz()
	data.ResetGames()
}

func SetTestURL(url *url.URL) {
	serverURL = url
}

func GetTestURL() url.URL {
	return *serverURL
}
