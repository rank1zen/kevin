-- Write your migrate up statements here

create table LiveMatch (
    match_id riot_match_id primary key,
    region riot_platform not null,
    date timestamp with time zone not null
);

create table LiveParticipant (
    match_id riot_match_id not null,
    puuid riot_puuid not null,
    primary key (match_id, puuid),

    champion_id int not null,
    runes int[11] not null,
    summoners int[2] not null,
    team_id int not null
);

create table LiveMatchStatus (
    match_id riot_match_id primary key,
    region riot_platform not null,
    date timestamp with time zone not null,
    expired boolean not null
);

---- create above / drop below ----

drop table LiveMatchStatus;
drop table LiveParticipant;
drop table LiveMatch;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
