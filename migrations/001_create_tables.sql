-- Write your migrate up statements here

create domain riot_puuid as char(78);

create domain riot_match_id as varchar(60);

create type riot_platform as enum (
        'BR1',
        'EUN1',
        'EUW1',
        'JP1',
        'KR',
        'LA1',
        'LA2',
        'NA1',
        'OC1',
        'TR1',
        'RU',
        'PH2',
        'SG2',
        'TH2',
        'TW2',
        'VN2'
);

create type league_tier as enum (
        'Iron',
        'Bronze',
        'Silver',
        'Gold',
        'Platinum',
        'Emerald',
        'Diamond',
        'Master',
        'Grandmaster',
        'Challenger'
);

create type league_division as enum (
        'IV',
        'III',
        'II',
        'I'
);

create table Summoner (
        puuid       riot_puuid primary key,
        name        varchar(32) not null,
        tagline     varchar(10) not null
);

create table RankStatus (
        rank_status_id serial primary key,
        puuid          riot_puuid not null references Summoner(puuid),
        effective_date timestamp without time zone not null,
        end_date       timestamp without time zone not null,
        is_current     boolean not null,
        is_ranked      boolean not null
);

create table RankDetail (
        rank_status_id int primary key references RankStatus(rank_status_id),
        wins           int not null,
        losses         int not null,
        tier           league_tier not null,
        division       league_division not null,
        lp             int not null
);

create table Match (
        match_id riot_match_id primary key,
        date     timestamp without time zone not null,
        duration interval not null,
        version  text not null,
        winner   int not null
);

create table Participant (
        match_id riot_match_id not null,
        puuid riot_puuid not null,
        primary key (match_id, puuid),

        team                   int     not null,
        champion               int     not null,
        champion_level         int     not null,
        summoners              int[2]  not null,
        runes                  int[11] not null,
        items                  int[7]  not null,
        kills                  int     not null,
        deaths                 int     not null,
        assists                int     not null,
        kill_participation     decimal not null,
        creep_score            int     not null,
        creep_score_per_minute decimal not null,
        damage_dealt           int     not null,
        damage_taken           int     not null,
        damage_delta_enemy     int     not null,
        damage_percentage_team decimal not null,
        gold_earned            int     not null,
        gold_delta_enemy       int     not null,
        gold_percentage_team   decimal not null,
        vision_score           int     not null,
        pink_wards_bought      int     not null
);

create type item_event_type as enum (
        'purchased',
        'sold'
);

create table ItemEvent (
        match_id          riot_match_id primary key,
        puuid             riot_puuid not null,
        in_game_timestamp interval not null,
        item_id           int not null,
        type              item_event_type not null
);

create table SkillEvent (
        match_id          riot_match_id primary key,
        puuid             riot_puuid not null,
        in_game_timestamp interval not null,
        spell_slot        int not null
);

---- create above / drop below ----

drop table Participant;
drop table Match;
drop table Matchlist;
drop table RankDetail;
drop table RankStatus;
drop table Summoner;
drop table ItemEvent;
drop table SkillEvent;

drop type league_division;
drop type league_tier;
drop type riot_platform;
drop type riot_puuid;
drop type riot_match_id;
drop type item_event_type;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
