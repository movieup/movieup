
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE mst_movie
(
  movie_id INT NOT NULL,
  title_jp VARCHAR(512) NOT NULL,
  title_en VARCHAR(512) NOT NULL,
  description VARCHAR(1024) NOT NULL,
  recommended_point INT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME,
  deleted_at DATETIME,
  PRIMARY KEY (movie_id)
);
CREATE TABLE tbl_scene
(
  movie_id INT NOT NULL,
  scene_id INT NOT NULL,
  recommended_number INT NOT NULL,
  description VARCHAR(1024) NOT NULL,
  memo VARCHAR(1024),
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME,
  deleted_at DATETIME,
  PRIMARY KEY (movie_id, scene_id),
  FOREIGN KEY (movie_id) REFERENCES mst_movie(movie_id)
);
CREATE TABLE tbl_scene_detail
(
  movie_id INT NOT NULL,
  scene_id INT NOT NULL,
  scene_detail_id INT NOT NULL,
  order_number INT NOT NULL,
  characterName VARCHAR(128),
  quote VARCHAR(1024),
  translated VARCHAR(1024),
  action VARCHAR(1024),
  PRIMARY KEY (movie_id, scene_id, scene_detail_id),
  FOREIGN KEY (movie_id, scene_id) REFERENCES tbl_scene(movie_id, scene_id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tbl_scene;
DROP TABLE mst_movie;
DROP TABLE tbl_scene_detail;
