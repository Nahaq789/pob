pob/
├── go.work
├── go.work.sum
├── docker-compose.yml
├── scripts/
│   └── init-db.sql
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
│   │   └── jwt.go
│   └── stats/
│       └── stats.go                # 実数値計算 (nature-adjusted)
│
├── dex/                            # 3層構造
│   ├── go.mod                      # module pob/dex
│   ├── Dockerfile
│   ├── cmd/
│   │   └── main.go
│   ├── gen/                        # dex.proto から生成（サーバ用）
│   │   ├── dex.pb.go
│   │   └── dex_grpc.pb.go
│   ├── internal/
│   │   ├── handler/                # Echo ハンドラ + gRPC サーバ実装
│   │   ├── service/
│   │   └── repository/
│   └── shared/
│       └── model/                  # DB モデル
│
├── user/                           # 3層構造
│   ├── go.mod                      # module pob/user
│   ├── Dockerfile
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   └── shared/
│       └── model/                  # DB モデル
│
├── box/                            # 3層構造
│   ├── go.mod                      # module pob/box
│   ├── Dockerfile
│   ├── cmd/
│   │   └── main.go
│   ├── gen/                        # box.proto（サーバ用）+ dex.proto（クライアント用）から生成
│   │   ├── box.pb.go
│   │   ├── box_grpc.pb.go
│   │   ├── dex.pb.go
│   │   └── dex_grpc.pb.go
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   └── shared/
│       └── model/                  # DB モデル
│
├── battle/                         # DDD 構造
│   ├── go.mod                      # module pob/battle
│   ├── Dockerfile
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── gen/                    # dex.proto + box.proto から生成（クライアント用）
│       │   ├── dex.pb.go
│       │   ├── dex_grpc.pb.go
│       │   ├── box.pb.go
│       │   └── box_grpc.pb.go
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
└── sync/                           # PokeAPI データ同期バッチ
    ├── go.mod                      # module pob/sync
    ├── Dockerfile
    └── cmd/
        └── main.go
