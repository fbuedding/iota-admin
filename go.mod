module github.com/fbuedding/iota-admin

go 1.21.1

require (
	github.com/Netflix/go-env v0.0.0-20220526054621-78278af1949d
	github.com/a-h/templ v0.2.639
	github.com/go-chi/chi/v5 v5.0.10
	github.com/google/uuid v1.4.0
	github.com/gorilla/schema v1.2.1
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/rs/zerolog v1.32.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/monoculum/formam v3.5.5+incompatible // indirect
	github.com/niemeyer/golang v0.0.0-20110826170342-f8c0f811cb19 // indirect
)

require (
	github.com/fbuedding/fiware-iot-agent-sdk v1.1.8
	github.com/go-chi/httprate v0.9.0
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/monoculum/formam/v3 v3.6.0
	golang.org/x/sys v0.16.0 // indirect
)

//wrongly published v1.4.0
retract v1.4.0
