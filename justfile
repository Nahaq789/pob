set dotenv-load := true

@run_user:
    cd ./user && go run ./...

@run_dex:
    cd ./dex && go run ./...
