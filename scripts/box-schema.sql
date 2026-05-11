\c box_db

CREATE TABLE boxes (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL,
    name       VARCHAR(50) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE box_pokemon (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    box_id          UUID        NOT NULL REFERENCES boxes(id) ON DELETE CASCADE,
    pokemon_id      INT         NOT NULL,  -- dex_db の pokemon.id
    nickname        VARCHAR(50),
    ability_id      INT         NOT NULL,  -- dex_db の abilities.id
    nature          VARCHAR(20) NOT NULL,
    held_item_id    INT,                   -- NULL: 道具なし
    -- 個体値
    iv_hp           INT         NOT NULL DEFAULT 31,
    iv_attack       INT         NOT NULL DEFAULT 31,
    iv_defense      INT         NOT NULL DEFAULT 31,
    iv_sp_attack    INT         NOT NULL DEFAULT 31,
    iv_sp_defense   INT         NOT NULL DEFAULT 31,
    iv_speed        INT         NOT NULL DEFAULT 31,
    -- 努力値
    ev_hp           INT         NOT NULL DEFAULT 0,
    ev_attack       INT         NOT NULL DEFAULT 0,
    ev_defense      INT         NOT NULL DEFAULT 0,
    ev_sp_attack    INT         NOT NULL DEFAULT 0,
    ev_sp_defense   INT         NOT NULL DEFAULT 0,
    ev_speed        INT         NOT NULL DEFAULT 0,
    -- 技スロット（NULL許容）
    move1_id        INT,
    move2_id        INT,
    move3_id        INT,
    move4_id        INT,
    created_at      TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE parties (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL,
    name       VARCHAR(50) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- パーティ内のポケモン（スロット1〜6）
CREATE TABLE party_pokemon (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    party_id        UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
    box_pokemon_id  UUID NOT NULL REFERENCES box_pokemon(id),
    slot            INT  NOT NULL CHECK (slot BETWEEN 1 AND 6),
    UNIQUE (party_id, slot)
);