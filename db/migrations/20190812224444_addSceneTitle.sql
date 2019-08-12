
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table tbl_scene add title varchar(512) after recommended_number;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table tbl_scene drop title;
