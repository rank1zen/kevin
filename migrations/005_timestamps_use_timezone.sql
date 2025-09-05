-- Write your migrate up statements here

alter table RankStatus
        alter column effective_date type timestamp with time zone using effective_date at time zone 'UTC';

alter table Match
        alter column date type timestamp with time zone using date at time zone 'UTC';

---- create above / drop below ----

alter table RankStatus
        alter column effective_date type timestamp without time zone using effective_date at time zone 'UTC';

alter table Match
        alter column date type timestamp without time zone using date at time zone 'UTC';

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
