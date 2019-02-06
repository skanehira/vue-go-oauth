# vue-go-oauth
Authentication Twitter using vue.js and Go sample

# requirements
- Go 1.11.5
- Docker 18.09.1
- Node.js 11.7.0

# DB
```shell
$ cd docker
$ docker-compose up -d
```

# Front
```shell
$ cd front
$ npm insatll
$ npm run build
```

# API Server
```shell
$ cd api
$ go build
$ ./api
```

## config
```yaml
appname: test
port: 80
dblog: false
applog: false
db:
    name:     test
    user:     test
    password: test
    port:     3306
    host:     127.0.0.1
twitter:
    token: "your app token"
    secret: "your app secret"
    requesturi: "https://api.twitter.com/oauth/request_token"
    authorizationuri: "https://api.twitter.com/oauth/authorize"
    tokenrequesturi: "https://api.twitter.com/oauth/access_token"
    callbackuri:  "http://localhost/twitter/callback"
```
