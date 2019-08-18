
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table mst_movie add amazon_affiliate_url varchar(1024) after recommended_point;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table mst_movie drop amazon_affiliate_url;
