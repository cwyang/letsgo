package main

import (
	"flag"
	"log"
	"net/http"
)

const html = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
  </head>
  <body>
    <h1> CORS demo </h1>
    <div id="output"></div>
    <script>
      document.addEventListener('DOMContentLoaded', function() {
        fetch("http://localhost:4040/v1/healthcheck").then(
          function (response) {
            response.text().then(function (text) {
              document.getElementById("output").innerHTML = text;
            });
          },
          function (err) {
              document.getElementById("output").innerHTML = err;
          }
        );
      });
    </script>
  </body>
</html>`

func main() {
	addr := flag.String("addr", ":9000", "server address")
	flag.Parse()
	log.Printf("starting server on %s", *addr)

	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	log.Fatal(err)
}
