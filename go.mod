module github.com/schulterklopfer/cyphernode_admin

go 1.12

// replace github.com/schulterklopfer/cyphernode_fauth => ../cyphernode_fauth

require (
	github.com/SatoshiPortal/cam v0.0.0-20210219205004-f45be2385b55
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/docker v20.10.3+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/go-resty/resty/v2 v2.5.0
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/schulterklopfer/cyphernode_fauth v0.0.0-20210702170312-3bdcc12bc919
	github.com/sirupsen/logrus v1.8.0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/text v0.3.5
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
	google.golang.org/grpc v1.39.0 // indirect
	gopkg.in/validator.v2 v2.0.0-20200605151824-2b28d334fa05
)
