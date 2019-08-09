
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE mst_dictionary
(
  dictionary_id INT NOT NULL,
  name VARCHAR(254),
  description VARCHAR(1024) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME,
  deleted_at DATETIME,
  PRIMARY KEY (dictionary_id)
);
CREATE TABLE tbl_scene_dictionary
(
  movie_id INT NOT NULL,
  scene_id INT NOT NULL,
  dictionary_id INT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME,
  deleted_at DATETIME,
  PRIMARY KEY (movie_id, scene_id, dictionary_id),
  FOREIGN KEY (movie_id, scene_id) REFERENCES tbl_scene(movie_id, scene_id),
  FOREIGN KEY (dictionary_id) REFERENCES mst_dictionary(dictionary_id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tbl_scene_dictionary;
DROP TABLE mst_dictionary;
