module github.com/fbuedding/iota-admin

go 1.21.1

require (
	github.com/Netflix/go-env v0.0.0-20220526054621-78278af1949d
	github.com/a-h/templ v0.2.432
	github.com/go-chi/chi/v5 v5.0.10
	github.com/google/uuid v1.4.0
	github.com/gorilla/schema v1.2.0
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/niemeyer/golang v0.0.0-20110826170342-f8c0f811cb19
	github.com/rs/zerolog v1.31.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

//wrongly published v1.4.0
retract v1.4.0
