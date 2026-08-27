-- Write your migrate up statements here

create extension if not exists pg_trgm;

create domain riot_puuid as char(78);

create table Summoner
(
    id              uuid        not null primary key default uuidv7(),
    puuid           riot_puuid  not null,
    region          text        not null,
    name            text        not null,
    tagline         text        not null,
    summoner_level  int         not null,
    profile_icon_id text        not null,
    last_updated    timestamptz not null,

    unique (puuid),
    unique (region, name, tagline)
);

create table Rank
(
    id            uuid        not null primary key default uuidv7(),
    puuid         riot_puuid  not null,
    wins          int         not null,
    losses        int         not null,
    tier          text        not null,
    division      text        not null,
    league_points int         not null,
    last_updated  timestamptz not null,

    unique (puuid)
);

create table RankHistory
(
    id            uuid        not null primary key default uuidv7(),
    puuid         riot_puuid  not null,
    valid_from    timestamptz not null,
    valid_to      timestamptz,
    wins          int         not null,
    losses        int         not null,
    tier          text        not null,
    division      text        not null,
    league_points int         not null
);

create table Match
(
    id        uuid        not null primary key default uuidv7(),
    region    text        not null,
    match_id  text        not null,
    version   text        not null,
    date      timestamptz not null,
    duration  interval    not null,
    winner_id text        not null,

    unique (match_id)
);

create table Participant
(
    id                uuid       not null primary key default uuidv7(),
    match_id          text       not null,
    puuid             riot_puuid not null,
    team_id           text       not null,
    team_position     text       not null,
    champion_id       text       not null,
    champion_level    int        not null,
    summoner_ids      text[2]    not null,
    rune_ids          text[11]   not null,
    item_ids          text[7]    not null,
    kills             int        not null,
    deaths            int        not null,
    assists           int        not null,
    creep_score       int        not null,
    damage_dealt      int        not null,
    damage_taken      int        not null,
    gold_earned       int        not null,
    vision_score      int        not null,
    pink_wards_bought int        not null,

    unique (match_id, puuid)
);

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
