# okotron
a cross chain telegram trading bot

### Pre requisites
- go 1.21+

### Build
```
make build
```

### test
```
make test
```

### Configure
set the below environment variables
```
TELEGRAM_BOT_TOKEN=<your-secret-token>
ENABLE_DEBUG_LOGS=true
GOOGLE_CLIENT_ID=<client-id>
GOOGLE_CLIENT_SECRET=<secret>
OKTO_CLIENT_API_KEY: <your-okto-client-key>,
REDIS_ADDR : <redis-address>,
REDIS_CMD_EXPIRY_IN_SEC: <command-expiry>,
CMC_KEY : <coin-market-cap-api-key>,
```

### Run
```
make run
```

### Deployment
- cross compile the binary to target machine. otionally, mordenize the application.
    - for example, makefile task compiles to a linux compatible machine
- run the binary in the production machine