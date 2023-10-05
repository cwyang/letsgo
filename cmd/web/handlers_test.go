package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestShowNote(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name string
		url  string
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

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	t.Log(csrfToken) // go test -v showes test log

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		want         int
		body         []byte
	}{
		{"valid sub", "Alice", "alice@mail.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "alice@mail.com", "validPa$$word", csrfToken, http.StatusOK, []byte("must be present")},
		{"Empty email", "Alice", "", "validPa$$word", csrfToken, http.StatusOK, []byte("must be present")},
		{"Empty password", "Alice", "alice@mail.com", "", csrfToken, http.StatusOK, []byte("must be present")},
		{"short password", "Alice", "alice@mail.com", "password", csrfToken, http.StatusOK, []byte("short")},	
		{"Dup email", "Alice", "dup@mail.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Duplicate")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			if code != tt.want {
				t.Errorf("want %d; got %d", tt.want, code)
			}
			if !bytes.Contains(body, tt.body) {
				t.Errorf("want body to contain %q", tt.body)
			}
		})
	}
}

func TestCreateNoteForm(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/note/create")
		if code != http.StatusSeeOther {
			t.Errorf("want %d; got %d", http.StatusSeeOther, code)
		}
		if headers.Get("Location") != "/user/login" {
			t.Errorf("want %s; got %s", "/user/login", headers.Get("Location"))
		}
	})

	t.Run("authenticated", func(t *testing.T) {
		_, _, body2 := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body2)
		form := url.Values{}
		form.Add("email", "kim@mail.com")
		form.Add("password", "")
		form.Add("csrf_token", csrfToken)

		ts.postForm(t, "/user/login", form)
		
		_, _, body := ts.get(t, "/note/create")
		target := []byte("<form action='/note/create' method='POST'>")
		if !bytes.Contains(body, target) {
			t.Errorf("want body to contain %q", target)
		}
	})
}
