module github.com/schulterklopfer/cyphernode_admin

go 1.12

replace github.com/SatoshiPortal/cam => /Users/jash/go/src/github.com/SatoshiPortal/cam

require (
	github.com/SatoshiPortal/cam v0.0.0-20200807091734-ccc5b29959eb
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.1
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-resty/resty/v2 v2.3.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.1.4-0.20181208214519-12bd4761fc66
	github.com/jinzhu/gorm v1.9.10
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20191106202628-ed6320f186d4
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	gopkg.in/validator.v2 v2.0.0-20191008145730-5614e8810ea7
)
