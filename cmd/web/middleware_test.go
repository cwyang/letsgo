package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	
	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	checkHeader := func(hn string, want string) {
		v := rs.Header.Get(hn)
		if v != want {
			t.Errorf("want %q; got %q", want, v)
		}
	}
	checkHeader("X-Frame-Options", "deny")
	checkHeader("X-XSS-Protection", "1; mode=block")
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "Pong")
	}
}
