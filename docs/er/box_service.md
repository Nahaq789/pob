erDiagram
    boxes {
        UUID id PK "ボックスID"
        UUID user_id "ユーザーID"
        VARCHAR_50 name "ボックス名"
        TIMESTAMPTZ created_at "作成日時"
        TIMESTAMPTZ updated_at "更新日時"
    }

    box_pokemon {
        UUID id PK "個体ID"
        UUID box_id FK "ボックスID"
        INT pokemon_id "ポケモンID"
        VARCHAR_50 nickname "ニックネーム"
        INT nature_id "性格ID"
        INT ability_id "特性ID"
        INT move1_id "技1ID"
        INT move2_id "技2ID"
        INT move3_id "技3ID"
        INT move4_id "技4ID"
        SMALLINT iv_hp "個体値HP"
        SMALLINT iv_attack "個体値こうげき"
        SMALLINT iv_defense "個体値ぼうぎょ"
        SMALLINT iv_sp_attack "個体値とくこう"
        SMALLINT iv_sp_defense "個体値とくぼう"
        SMALLINT iv_speed "個体値すばやさ"
        SMALLINT ev_hp "努力値HP"
        SMALLINT ev_attack "努力値こうげき"
        SMALLINT ev_defense "努力値ぼうぎょ"
        SMALLINT ev_sp_attack "努力値とくこう"
        SMALLINT ev_sp_defense "努力値とくぼう"
        SMALLINT ev_speed "努力値すばやさ"
        TIMESTAMPTZ created_at "作成日時"
        TIMESTAMPTZ updated_at "更新日時"
    }

    parties {
        UUID id PK "パーティID"
        UUID user_id "ユーザーID"
        VARCHAR_50 name "パーティ名"
        TIMESTAMPTZ created_at "作成日時"
        TIMESTAMPTZ updated_at "更新日時"
    }

    party_pokemon {
        UUID id PK "パーティポケモンID"
        UUID party_id FK "パーティID"
        UUID box_pokemon_id FK "個体ID"
        SMALLINT slot "スロット"
    }

    boxes ||--o{ box_pokemon : "has"
    parties ||--o{ party_pokemon : "has"
    box_pokemon ||--o{ party_pokemon : "has"
