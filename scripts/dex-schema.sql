\c dex_db

CREATE TABLE types (
    id   INT         PRIMARY KEY,  -- PokeAPI ID
    name VARCHAR(50) NOT NULL
);

CREATE TABLE abilities (
    id   INT          PRIMARY KEY,  -- PokeAPI ID
    name VARCHAR(100) NOT NULL
);

CREATE TABLE pokemon (
    id              INT          PRIMARY KEY,  -- PokeAPI ID
    name            VARCHAR(100) NOT NULL,
    type1_id        INT          NOT NULL REFERENCES types(id),
    type2_id        INT          REFERENCES types(id),
    base_hp         INT          NOT NULL,
    base_attack     INT          NOT NULL,
    base_defense    INT          NOT NULL,
    base_sp_attack  INT          NOT NULL,
    base_sp_defense INT          NOT NULL,
    base_speed      INT          NOT NULL
);

CREATE TABLE pokemon_abilities (
    pokemon_id INT     NOT NULL REFERENCES pokemon(id),
    ability_id INT     NOT NULL REFERENCES abilities(id),
    slot       INT     NOT NULL,
    is_hidden  BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (pokemon_id, ability_id)
);

CREATE TABLE moves (
    id           INT          PRIMARY KEY,  -- PokeAPI ID
    name         VARCHAR(100) NOT NULL,
    type_id      INT          NOT NULL REFERENCES types(id),
    damage_class VARCHAR(20)  NOT NULL,     -- physical / special / status
    power        INT,                       -- NULL: 変動・固定ダメージ技など
    accuracy     INT,                       -- NULL: 必中技
    pp           INT          NOT NULL
);

CREATE TABLE pokemon_moves (
    pokemon_id INT NOT NULL REFERENCES pokemon(id),
    move_id    INT NOT NULL REFERENCES moves(id),
    PRIMARY KEY (pokemon_id, move_id)
);