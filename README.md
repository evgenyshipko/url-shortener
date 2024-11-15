### Golang Url Shortener

#### Get started:
1. create .env file with contents like in .env.sample
1. `go mod tidy` - install modules
2.  `go build url-shortener/cmd` - run server

#### Features

- save url by alias POST /url
- redirect by alias GET /r/{alias}

#### Future plans:
- transfer from sqlite to postgres (with migrations, seeds, e.t.c)
- write more functional test
- implement Delete api endpoint
- implement jwt auth