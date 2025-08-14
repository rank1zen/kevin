-- Write your migrate up statements here

create type team_position as enum (
        'Top',
        'Jungle',
        'Middle',
        'Bottom',
        'Support'
);

alter table Participant
        add column position team_position not null;

---- create above / drop below ----

alter table Participant
        drop column position if exists;

drop type item_event_type;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
