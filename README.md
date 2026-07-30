# bookmark-service

bookmark-service/
│
├── cmd/
│   ├── api/
│   │   └── main.go              # base path, set up config, redis client, and engine
│   └── test/
│       └── main.go              # main test new feature
│
├── docs/
│   ├── docs.go                  # swagger docs
│   ├── swagger.json
│   └── swagger.yaml             
│
├── internal/                    # Business logic
│   ├── api/                 
│   │   ├── engine.go            # create a new engine (engine, redis client, config), start the application, and initializes the routes
│   │   └── config.go            # creates a new config for api (app port)
│   │
│   ├── handler/                 # HTTP handlers
│   │   ├── genpass.go           # using genpass service for handling api GET /genpass
│   │   ├── genpass_test.go      # unit test for genpass handler
│   │   └── shortenurl.go        # using shortenurl service for handling api POST /v1/links/shorten
│   │
│   ├── integration_test/                   
│   │   └── genpass_ep_test.go   # integration test for handling api GET /genpass
│   │
│   ├── model/                   
│   │   └── 
│   │
│   ├── repository/              # Database layer
│   │   ├── urlstorage.go        # using redis database to store the url
│   │   └── urlstorage_test.go   # unit test for url storage
│   │
│   └── service/                 # Business logic
│       ├── mocks/               
│       │   └── shortenurl.go    # mock test service genpass for genpass handler
│       ├── genpass.go           # generate pass using charset
│       ├── genpass_test.go      # unit test for generate pass
│       └── shortenurl.go        # generate pass for url and store in repository
│
├── pkg/                         # reused package
│   └── redis/
│       ├── client.go            # creates a new redis client
│       ├── config.go            # creates a new config for redis
│       └── mock.go              # creates a mock redis client
│
├── .env
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── README.md