package main

import (
    "./structs"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
    "html/template"
    "net/http"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "log"
    "os"
    "github.com/BurntSushi/toml"

    "fmt"
    "strings"
    "encoding/base64"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/gmail/v1"
)

// 設定ファイルへの実行ファイルからの相対パスを指定
// ここだけは環境によって変更すること
const confDir = "docker/golang/config/env/"

// HTMLタグをエスケープせずそのまま出力する
func getNoEscapedString(html string) template.HTML {
    return template.HTML(html)
}

func main() {
    // タイムゾーンを日本に設定
    jstLocation, error := time.LoadLocation("Asia/Tokyo")
    if error != nil {
        jstLocation = time.FixedZone("Asia/Tokyo", 9 * 60 * 60)
    }
    time.Local = jstLocation

    router := gin.Default()

    // 環境変数を取得
    envConf := getConfing()

    // 静的ファイルをマッピング
    router.Static("/resources", envConf.Etc.SrcPath + "resources")

    // viewテンプレートを読み込み
    router.LoadHTMLGlob(envConf.Etc.SrcPath + "templates/*")

    // トップページ（記事新着一覧）
    router.GET("/", func(context *gin.Context) {
        html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/index.tmpl.html"))
        router.SetHTMLTemplate(html)

        db := gormConnect()
        defer db.Close()

        // シーン一覧を取得
        scenes := []structs.Scene{}
        db.Preload("Movie").Limit(10).Order("created_at desc").Find(&scenes)
        // シーン詳細をタグ出力可能な形へ編集
        for index, scene := range scenes {
            scenes[index].DescriptionHtml = getNoEscapedString(scene.Description)
        }

        // タグ一覧を取得
        tags := []structs.Tag{}
        db.Find(&tags)

        context.HTML(http.StatusOK, "base.tmpl.html", gin.H{ "scenes": scenes, "fiveList": []int {1,2,3,4,5}, "tags": tags, })
    })

    // 記事取得API
    router.GET("/news/pages/:pageNo/", func(context *gin.Context) {
        // パラメータ取得
        pageNo, _ := strconv.Atoi(context.Param("pageNo"))

        db := gormConnect()
        defer db.Close()

        // シーン一覧を取得
        scenes := []structs.Scene{}
        db.Preload("Movie").Limit(2).Offset(2 * (pageNo - 1)).Order("created_at desc").Find(&scenes)
        // 作成日を画面表示用に整形
        copiedScenes := []structs.Scene{}
        for _, scene := range scenes {
            var copiedScene = scene
            copiedScene.CreatedAtString = scene.GetFormattedCreatedAt()
            copiedScenes = append(copiedScenes, copiedScene)
        }

        context.JSON(http.StatusOK, gin.H{ "scenes": copiedScenes, "fiveList": []int {1,2,3,4,5}, })
    })

    // タグに紐づく記事一覧画面
    router.GET("/tags/:tagId/", func(context *gin.Context) {
        // パラメータ取得
        tagId, _ := strconv.Atoi(context.Param("tagId"))

        html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/tagScenes.tmpl.html"))
        router.SetHTMLTemplate(html)

        db := gormConnect()
        defer db.Close()

        // タグに紐づくmovie_idリストを取得
        tag := structs.Tag{}
        db.Preload("SceneTags").Order("created_at desc").First(&tag, "tag_id=?", tagId)
        var sceneIds = []int{}
        for _, sceneTag := range tag.SceneTags {
            sceneIds = append(sceneIds, sceneTag.SceneId)
        }
        // movie_idを持つシーンを全て取得
        scenes := []structs.Scene{}
        db.Preload("Movie").Where("scene_id in (?)", sceneIds).Preload("Movie").Find(&scenes)
        // シーン詳細をタグ出力可能な形へ編集
        for index, scene := range scenes {
            scenes[index].DescriptionHtml = getNoEscapedString(scene.Description)
        }

        // タグ一覧を取得
        tags := []structs.Tag{}
        db.Find(&tags)

        context.HTML(http.StatusOK, "base.tmpl.html", gin.H{ "scenes": scenes, "fiveList": []int {1,2,3,4,5}, "tags": tags, "tag": tag })
    })

    // シーン詳細画面
    router.GET("/scenes/:sceneId/", func(context *gin.Context) {
        // パラメータ取得
        sceneId, _ := strconv.Atoi(context.Param("sceneId"))

        html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/scene.tmpl.html"))
        router.SetHTMLTemplate(html)

        db := gormConnect()
        defer db.Close()

        // シーンに紐づくタグ一覧取得
        sceneTags := []structs.SceneTag{}
        db.Preload("Tag").Find(&sceneTags, "scene_id=?", sceneId)

        // シーン詳細一覧取得
        sceneDetails := []structs.SceneDetail{}
        db.Order("order_number asc").Find(&sceneDetails, "scene_id=?", sceneId)

        // 映画情報取得
        movie := structs.Movie{}
        if len(sceneDetails) > 0 {
            db.Find(&movie, "movie_id=?", sceneDetails[0].MovieId)
        }

        // シーン情報取得
        scene := structs.Scene{}
        db.Find(&scene, "scene_id=?", sceneId)
        // シーン詳細をタグ出力可能な形へ編集
        scene.DescriptionHtml = getNoEscapedString(scene.Description)

        // キャスト情報取得
        sceneCasts := []structs.SceneCast{}
        db.Find(&sceneCasts, "scene_id=?", sceneId)

        // 辞書情報取得
        sceneDictionaries := []structs.SceneDictionary{}
        //db.Find(&sceneDictionaries, "scene_id=?", sceneId)
        //db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
        db.Debug().Table("tbl_scene_dictionary").Select("tbl_scene_dictionary.dictionary_id, mst_dictionary.name, mst_dictionary.description").Joins("inner join mst_dictionary on tbl_scene_dictionary.dictionary_id = mst_dictionary.dictionary_id").Where("scene_id = ?", sceneId).Order("mst_dictionary.name asc").Scan(&sceneDictionaries)

        // タグ一覧を取得
        tags := []structs.Tag{}
        db.Find(&tags)

        context.HTML(http.StatusOK, "base.tmpl.html", gin.H{
            "tags": tags,
            "sceneTags": sceneTags,
            "sceneDetails": sceneDetails,
            "movie": movie,
            "scene": scene,
            "sceneCasts": sceneCasts,
            "movieDescription": getNoEscapedString(movie.Description),
            "sceneDescription": getNoEscapedString(scene.Description),
            "sceneDictionaries": sceneDictionaries,
            "memo": getNoEscapedString(scene.Memo) })
    })
    // お問い合わせ
    router.GET("/inquiry/", func(context *gin.Context) {
        html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/inquiry.tmpl.html"))
        router.SetHTMLTemplate(html)

        // "/inquiry/send/"からのリダイレクト時にエラーがある場合はerrorItemにエラーの項目名が設定されてくる
        errorItem := context.DefaultQuery("errorItem", "")
        mailAddress := context.DefaultQuery("mailAddress", "")
        inquiry := context.DefaultQuery("inquiry", "")

        db := gormConnect()
        defer db.Close()

        // タグ一覧を取得
        tags := []structs.Tag{}
        db.Find(&tags)

        context.HTML(http.StatusOK, "base.tmpl.html", gin.H{ "tags": tags, "errorItem": errorItem, "mailAddress": mailAddress, "inquiry": inquiry })
    })
    // お問い合わせ（確認画面）
    router.POST("/inquiry/send/", func(context *gin.Context) {
        db := gormConnect()
        defer db.Close()

        // タグ一覧を取得
        tags := []structs.Tag{}
        db.Find(&tags)

        // メールアドレス取得
        mailAddress := context.PostForm("mailAddress")
        // 問い合わせ内容取得
        inquiry := context.PostForm("inquiry")

        // validation
        if len(mailAddress) < 3 || 64 < len(mailAddress) {
            // メールアドレスが3〜64バイトでない場合はエラー
            html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/inquiry.tmpl.html"))
            router.SetHTMLTemplate(html)

            context.Redirect(http.StatusMovedPermanently, "/inquiry/?errorItem=mailAddress&mailAddress=" + mailAddress + "&inquiry=" + inquiry)
            return
        }
        if len(inquiry) < 10 || 2048 < len(inquiry) {
            // お問い合わせが10〜2048バイトでない場合はエラー
            html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/inquiry.tmpl.html"))
            router.SetHTMLTemplate(html)

            context.Redirect(http.StatusMovedPermanently, "/inquiry/?errorItem=inquiry&mailAddress=" + mailAddress + "&inquiry=" + inquiry)
            return
        }

        // メール送信
        sendMail(mailAddress, inquiry)

        html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/inquiryDone.tmpl.html"))
        router.SetHTMLTemplate(html)

        context.HTML(http.StatusOK, "base.tmpl.html", gin.H{ "tags": tags, "mailAddress": mailAddress, "inquiry": inquiry })
    })
    // 当サイトについて
    router.GET("/about/site/", func(context *gin.Context) {
        html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/aboutSite.tmpl.html"))
        router.SetHTMLTemplate(html)

        db := gormConnect()
        defer db.Close()

        // タグ一覧を取得
        tags := []structs.Tag{}
        db.Find(&tags)

        context.HTML(http.StatusOK, "base.tmpl.html", gin.H{ "tags": tags })
    })
    // 当サイトについて（技術）
    router.GET("/about/technology/", func(context *gin.Context) {
        html := template.Must(template.ParseFiles(envConf.Etc.SrcPath + "templates/base.tmpl.html", envConf.Etc.SrcPath + "templates/aboutTec.tmpl.html"))
        router.SetHTMLTemplate(html)

        db := gormConnect()
        defer db.Close()

        // タグ一覧を取得
        tags := []structs.Tag{}
        db.Find(&tags)

        context.HTML(http.StatusOK, "base.tmpl.html", gin.H{ "tags": tags })
    })

    router.Run(":8080")
}

func gormConnect() *gorm.DB {
    // 環境変数を取得
    envConf := getConfing()
    DBMS     := "mysql"
    USER     := envConf.DB.User
    PASS     := envConf.DB.Password
    PROTOCOL := envConf.DB.Protocol
    DBNAME   := envConf.DB.DBname

    CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME
    db,err := gorm.Open(DBMS, CONNECT + "?parseTime=true")

    if err != nil {
        panic(err.Error())
    }
    return db
}

func sendMail(mailAddress string, inquiry string) {
    // 環境変数を取得
    envConf := getConfing()

    config := oauth2.Config{
        ClientID:     envConf.Gmail.ClientID,
        ClientSecret: envConf.Gmail.ClientSecret,
        Endpoint:     google.Endpoint,
        RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
        Scopes:       []string{"https://mail.google.com/"},
    }

    expiry, _ := time.Parse("2006-01-02", "2017-07-11")
    token := oauth2.Token{
        AccessToken:  envConf.Gmail.AccessToken,
        TokenType:    "Bearer",
        RefreshToken: envConf.Gmail.RefreshToken,
        Expiry:       expiry,
    }

    client := config.Client(oauth2.NoContext, &token)

    srv, err := gmail.New(client)
    if err != nil {
        log.Fatalf("Unable to retrieve gmail Client %v", err)
    }

    temp := []byte("From: 'me'\r\n" +
        "reply-to: " + mailAddress + "\r\n" +
        "To: " + envConf.Gmail.MailAddress + "\r\n" +
        "Subject: MOViE UP inquiry\r\n" +
        "\r\n" + "MailAddress=" + mailAddress + "\r\n==========================\r\n" + inquiry)

    var message gmail.Message
    message.Raw = base64.StdEncoding.EncodeToString(temp)
    message.Raw = strings.Replace(message.Raw, "/", "_", -1)
    message.Raw = strings.Replace(message.Raw, "+", "-", -1)
    message.Raw = strings.Replace(message.Raw, "=", "", -1)

    _, err = srv.Users.Messages.Send("me", &message).Do()
    if err != nil {
        fmt.Printf("%v", err)
    }
}

func getConfing() structs.Config {
    var err error

    appMode := os.Getenv("APP_MODE") // 動作環境をAPP_MODEという形で環境変数に格納する
    if appMode == "" {
        panic("failed to get application mode, check whether APP_MODE is set.")
    }

    // Get configuration
    conf, err := newConfig(confDir, appMode) // 引数に渡す
    if err != nil {
        panic(err.Error())
    }

    return conf
}

// NewConfig return configuration struct.
func newConfig(path string, appMode string) (structs.Config, error) {
    var conf structs.Config

    confPath := path + appMode + ".toml" // tomlファイルを読み設定情報を取得
    if _, err := toml.DecodeFile(confPath, &conf); err != nil {
        return conf, err
    }

    return conf, nil
}
