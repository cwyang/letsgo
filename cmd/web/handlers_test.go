package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// new httptest.ResponseRecorder
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "Pong" {
		t.Errorf("want body to equal %q", "Pong")
	}
}

func TestPing2(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "Pong" {
		t.Errorf("want body to equal %q", "Pong")
	}
}

func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close();

	tests := []struct {
		name string
		url string
		want int
		body []byte
	}{
		{"valid id", "/note/1", http.StatusOK, []byte("Quick Brown Fox Jumps")},
		{"no id", "/note/2", http.StatusNotFound, nil},
		{"negative id", "/note/-1", http.StatusNotFound, nil},
		{"decimal id", "/note/1.23", http.StatusNotFound, nil},
		{"string id", "/note/foo", http.StatusNotFound, nil},
		{"empty id", "/note/", http.StatusNotFound, nil},
		{"trailing slash", "/note/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.url)

			if code != tt.want {
				t.Errorf("want %d; got %d", tt.want, code)
			}
			if !bytes.Contains(body, tt.body) {
				t.Errorf("want body to contain %q", tt.body)
			}
		})
	}
}
