-- Write your migrate up statements here

alter table RankStatus
        drop constraint rankstatus_puuid_fkey;

---- create above / drop below ----

alter table RankStatus
        add constraint rankstatus_puuid_fkey foreign key (puuid) references Summoner (puuid);

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
