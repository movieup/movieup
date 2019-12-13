
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table tbl_scene add short_title varchar(256) after title;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table tbl_scene drop short_title;
