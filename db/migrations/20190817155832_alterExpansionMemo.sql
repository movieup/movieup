
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE tbl_scene MODIFY memo VARCHAR(5096);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE tbl_scene MODIFY memo VARCHAR(1024);
