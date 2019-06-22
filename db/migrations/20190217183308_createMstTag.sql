
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE mst_tag
(
  tag_id INT NOT NULL,
  name VARCHAR(256) NOT NULL,
  order_number INT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME,
  deleted_at DATETIME,
  PRIMARY KEY (tag_id)
);
CREATE TABLE mst_scene_tag
(
  movie_id INT NOT NULL,
  scene_id INT NOT NULL,
  tag_id INT NOT NULL,
  PRIMARY KEY (movie_id, scene_id, tag_id),
  FOREIGN KEY (movie_id, scene_id) REFERENCES tbl_scene(movie_id, scene_id),
  FOREIGN KEY (tag_id) REFERENCES mst_tag(tag_id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE mst_movie_tag;
DROP TABLE mst_tag;
