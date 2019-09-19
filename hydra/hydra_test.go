package hydra_test

import (
	"github.com/schulterklopfer/cyphernode_admin/globals"
	"github.com/schulterklopfer/cyphernode_admin/hydra"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)


/*
func TestHelloHandler(t *testing.T) {


	//res, err := client.Get(ts.URL)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "http://localhsot", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloHandler(client, "http://localhost"))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
*/

func TestGetConsentRequest(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.GetConsentRequest( client, "123456" )

}

func TestAcceptConsentRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.AcceptConsentRequest( client, "123456", nil )

}

func TestRejectConsentRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.RejectConsentRequest( client, "123456", nil )

}

func TestGetLoginRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.GetLoginRequest( client, "123456" )

}

func TestAcceptLoginRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.AcceptConsentRequest( client, "123456", nil )

}

func TestRejectLoginRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.RejectConsentRequest( client, "123456", nil )

}

func TestGetLogoutRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.GetLogoutRequest( client, "123456" )

}

func TestAcceptLogoutRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.AcceptLogoutRequest( client, "123456", nil )

}

func TestRejectLogoutRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check some request stuff here
		}))
	defer ts.Close()
	_ = os.Setenv( globals.HYDRA_ADMIN_URL_ENV_KEY, ts.URL )

	client := ts.Client()
	_,_ = hydra.RejectLogoutRequest( client, "123456", nil )

}