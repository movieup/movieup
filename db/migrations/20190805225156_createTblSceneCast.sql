
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE tbl_scene_cast
(
  movie_id INT NOT NULL,
  scene_id INT NOT NULL,
  order_number INT NOT NULL,
  character_name VARCHAR(254),
  cast_name VARCHAR(256),
  description VARCHAR(1024) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME,
  deleted_at DATETIME,
  PRIMARY KEY (movie_id, scene_id, order_number),
  FOREIGN KEY (movie_id, scene_id) REFERENCES tbl_scene(movie_id, scene_id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tbl_scene_cast;
