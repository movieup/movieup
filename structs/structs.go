package structs

import (
    "time"
    "html/template"
)

type Movie struct {
    //gorm.Model
    MovieId       int    `gorm:"column:movie_id"`
    TitleJp       string `gorm:"column:title_jp"`
    TitleEn       string `gorm:"column:title_en"`
    Description   string `gorm:"column:description"`
    Scenes        []Scene `gorm:"foreignkey:movie_id;association_foreignkey:movie_id"`
}
func (m *Movie) TableName() string {
	return "mst_movie"
}

type Scene struct {
    //gorm.Model
    MovieId             int         `gorm:"column:movie_id;primary_key"`
    Movie               Movie       `gorm:"foreignkey:movie_id;association_foreignkey:movie_id"`
    SceneId             int         `gorm:"column:scene_id"`
    RecommendedNumber   int         `gorm:"column:recommended_number"`
    Title               string      `gorm:"column:title"`
    Description         string      `gorm:"column:description"`
    DescriptionHtml     template.HTML
    Memo                string      `gorm:"column:memo"`
    CreatedAt           time.Time   `gorm:"column:created_at"`
    CreatedAtString     string
    FormattedCreatedAt  string
}
func (m *Scene) TableName() string {
	return "tbl_scene"
}
func (scene *Scene) GetFormattedCreatedAt() string {
    return scene.CreatedAt.Format("2006-01-02")
}

type Tag struct {
    TagId       int `gorm:"column:tag_id"`
    Name        string `gorm:"column:name"`
    OrderNumber int `gorm:"column:order_number"`
    SceneTags   []SceneTag `gorm:"foreignkey:tag_id;association_foreignkey:tag_id"`
}
func (m *Tag) TableName() string {
	return "mst_tag"
}

type SceneTag struct {
    MovieId int `gorm:"column:movie_id;primary_key"`
    SceneId int `gorm:"column:scene_id;primary_key"`
    TagId   int `gorm:"column:tag_id"`
    //Movies   []Movie `gorm:"foreignkey:movie_id;association_foreignkey:movie_id"`
    Scenes  []Scene `gorm:"foreignkey:movie_id,scene_id;association_foreignkey:movie_id,scene_id"`
    Tag     Tag `gorm:"foreignkey:tag_id;association_foreignkey:tag_id"`
}
func (m *SceneTag) TableName() string {
	return "mst_scene_tag"
}
type SceneDetail struct {
    MovieId int `gorm:"column:movie_id;primary_key"`
    SceneId int `gorm:"column:scene_id;primary_key"`
    SceneDetailId int `gorm:"column:scene_detail_id;primary_key"`
    OrderNumber int `gorm:"column:order_number"`
    CharacterName    string      `gorm:"column:characterName"`
    Quote    string      `gorm:"column:quote"`
    Translated    string      `gorm:"column:translated"`
    Action    string      `gorm:"column:action"`
}
func (m *SceneDetail) TableName() string {
	return "tbl_scene_detail"
}
type SceneCast struct {
    MovieId int `gorm:"column:movie_id;primary_key"`
    SceneId int `gorm:"column:scene_id;primary_key"`
    OrderNumber int `gorm:"column:order_number;primary_key"`
    CharacterName    string      `gorm:"column:character_name"`
    CastName    string      `gorm:"column:cast_name"`
    Description    string      `gorm:"column:description"`
}
func (m *SceneCast) TableName() string {
	return "tbl_scene_cast"
}
type SceneDictionary struct {
    MovieId int `gorm:"column:movie_id;primary_key"`
    SceneId int `gorm:"column:scene_id;primary_key"`
    DictionaryId int `gorm:"column:dictionary_id;primary_key"`
    //Dictionary  Dictionary `gorm:"foreignkey:dictionary_id;association_foreignkey:dictionary_id"`
    // mst_dictionaryをjoinでひっぱってくるためName、Descriptionをここへ定義
    Name string
    Description    string
}
func (m *SceneDictionary) TableName() string {
	return "tbl_scene_dictionary"
}
