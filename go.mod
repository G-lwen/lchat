module lchat

go 1.14

require (
	github.com/appleboy/gin-jwt/v2 v2.6.3
	github.com/boj/redistore v0.0.0-20180917114910-cd5dcc76aeff // indirect
	github.com/gin-gonic/contrib v0.0.0-20191209060500-d6e26eeaa607
	github.com/gin-gonic/gin v1.6.3
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/russross/blackfriday.v2 v2.0.1
	gopkg.in/yaml.v2 v2.3.0
)

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
