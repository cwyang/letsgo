# letsgo
Let's Go(lang) playground

## Golang setup

```
$ mkdir work; cd work
$ touch main.go
$ go mod init my.url/my_module_name

```

## Run
```
$ go run main.go               # file
$ go run .                     # path
$ go run my.url/my_module_name # module
```

## project structure
* cmd: application specific code
* pkg: non application specific code, like validation helpers and SQL model.
* ui: user-interface asset (non-golang)
