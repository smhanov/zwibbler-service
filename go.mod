module zwibbler.com/zwibbler

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.2
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.3
	github.com/gorilla/websocket v1.5.0
	github.com/kardianos/service v1.1.0
	github.com/lib/pq v1.10.9 // indirect
	github.com/smhanov/zwibserve v0.0.0-20240502203354-7846f06df9f2
)

//replace github.com/smhanov/zwibserve v0.0.0-20231011211758-2b6f0e54bcc5 => ./zwibserve
