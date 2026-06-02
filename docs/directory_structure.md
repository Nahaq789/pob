pob/
├── go.work
├── go.work.sum
├── docker-compose.yml
├── justfile
├── scripts/
│   ├── dex-schema.sql
│   ├── user-schema.sql
│   ├── box-schema.sql
│   ├── setup.sh
│   ├── truncate.sh
│   └── test-dex.sh
│
├── proto/                          # 契約定義のみ（言語非依存）
│   ├── dex.proto
│   └── box.proto
│
├── pkg/                            # サービス横断の共有コード
│   ├── go.mod                      # module pob/pkg
│   ├── logger/
│   │   └── logger.go
│   ├── jwt/
│   │   ├── jwt.go                  # VerifyToken, ExtractBearerToken
│   │   └── claims.go
│   ├── tracing/
│   │   └── context_key.go
│   ├── interceptor/
│   │   └── hmac/
│   │       ├── client.go           # HmacClientInterceptor
│   │       └── server.go           # HmacServerInterceptor
│   ├── redis/
│   │   └── redis_client.go
│   └── stats/
│       └── stats.go                # 実数値計算（nature補正済み）
│
├── dex/                            # gRPCのみ（:9091）
│   ├── go.mod                      # module pob/dex
│   ├── Dockerfile
│   ├── cmd/
│   │   ├── main.go
│   │   └── test/
│   │       └── main.go             # gRPC動作確認用
│   ├── proto/                      # dex.proto から生成（サーバ用）
│   │   ├── dex.pb.go
│   │   └── dex_grpc.pb.go
│   └── internal/
│       ├── handler/
│       │   └── dex_handler.go      # gRPCサーバ実装
│       ├── service/
│       │   ├── pokemon_service.go
│       │   ├── move_service.go
│       │   ├── ability_service.go
│       │   └── item_service.go
│       ├── repository/
│       │   ├── pokemon_repository.go
│       │   ├── move_repository.go
│       │   ├── ability_repository.go
│       │   └── item_repository.go
│       ├── model/
│       │   ├── pokemon.go
│       │   ├── move.go
│       │   ├── ability.go
│       │   ├── item.go
│       │   └── entity/             # GORMエンティティ
│       │       ├── pokemon.go
│       │       ├── move.go
│       │       ├── ability.go
│       │       ├── item.go
│       │       ├── nature.go
│       │       ├── type.go
│       │       ├── pokemon_ability.go
│       │       └── pokemon_move.go
│       └── shared/
│           └── db_client.go        # pgxpool + GORM
│
├── user/                           # REST（:8082）
│   ├── go.mod                      # module pob/user
│   ├── Dockerfile
│   ├── pem/                        # RSA鍵ペア（.gitignore除外）
│   ├── cmd/
│   │   ├── main.go
│   │   └── router.go
│   └── internal/
│       ├── handler/
│       │   ├── user_handler.go
│       │   └── auth_handler.go
│       ├── service/
│       │   ├── user_service.go
│       │   ├── auth_service.go
│       │   └── dto/
│       │       ├── user/
│       │       └── auth/
│       ├── repository/
│       │   ├── user_repository.go
│       │   ├── auth_repository.go
│       │   ├── refresh_token_repository.go
│       │   ├── hash.go
│       │   └── port/
│       │       └── transaction_manager.go
│       ├── model/
│       │   ├── user.go
│       │   ├── auth.go
│       │   ├── jwt.go
│       │   ├── refresh_token.go
│       │   └── apperror/
│       ├── middleware/
│       │   └── trace_middleware.go
│       └── shared/
│           ├── db_client.go        # pgxpool
│           └── transaction.go
│
├── box/                            # REST（:8083）+ gRPCサーバ（:9093）
│   ├── go.mod                      # module pob/box
│   ├── Dockerfile
│   ├── cmd/
│   │   ├── main.go
│   │   ├── router.go
│   │   └── di/                     # google/wire によるDI
│   │       ├── wire.go
│   │       ├── wire_gen.go
│   │       ├── box.go
│   │       ├── box_pokemon.go
│   │       ├── party.go
│   │       └── party_pokemon.go
│   ├── proto/                      # box.proto（サーバ用）+ dex.proto（クライアント用）から生成
│   │   ├── box.pb.go
│   │   ├── box_grpc.pb.go
│   │   ├── dex.pb.go
│   │   └── dex_grpc.pb.go
│   └── internal/
│       ├── handler/
│       │   ├── box_handler.go
│       │   ├── party_handler.go
│       │   ├── dex_handler.go      # dex gRPCクライアント経由のREST proxy
│       │   ├── grpc_handler.go     # BoxService gRPCサーバ実装（GetParty）
│       │   └── helper.go
│       ├── service/
│       │   ├── box_service.go
│       │   ├── box_pokemon_service.go
│       │   ├── party_service.go
│       │   ├── party_pokemon_service.go
│       │   └── dto/
│       │       ├── box.go
│       │       ├── box_pokemon.go
│       │       └── party.go
│       ├── repository/
│       │   ├── box_repository.go
│       │   ├── box_pokemon_repository.go
│       │   ├── party_repository.go
│       │   └── party_pokemon_repository.go
│       ├── model/
│       │   ├── box.go
│       │   ├── box_pokemon.go
│       │   ├── party.go
│       │   └── party_pokemon.go
│       ├── middleware/
│       │   ├── auth_middleware.go
│       │   └── trace_middleware.go
│       └── shared/
│           └── db_client.go        # pgxpool
│
├── battle/                         # REST+WebSocket（:8084）/ DDD構成
│   ├── go.mod                      # module pob/battle
│   ├── Dockerfile
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── gen/                    # dex.proto + box.proto から生成（クライアント用）
│       ├── presentation/
│       │   └── websocket/
│       ├── application/
│       ├── domain/
│       │   ├── battle/
│       │   ├── damage/
│       │   └── move/
│       └── infrastructure/
│           └── redis/
│
└── sync/                           # PokeAPI データ同期バッチ（Cobra CLI）
    ├── go.mod                      # module pob/sync
    ├── Dockerfile
    ├── cmd/
    │   └── main.go                 # csv / sync コマンド
    ├── client.go
    ├── csv.go
    ├── p.go                        # 並列fetch（errgroup）
    ├── range.go
    ├── util.go
    ├── type.go
    ├── item.go
    ├── ability.go
    ├── move.go
    ├── pokemon.go
    ├── pokemon_ability.go
    └── pokemon_move.go
