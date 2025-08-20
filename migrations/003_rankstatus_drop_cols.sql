-- Write your migrate up statements here

alter table RankStatus
        drop column end_date,
        drop column is_current;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
