package main

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"time"
	"log"

	"github.com/labstack/echo/v4"
)

//go:embed static
var staticFiles embed.FS

//go:embed templates
var templates embed.FS

//--- -s-
type Comment struct {
	Content string    `json:"content"`
	Created time.Time `json:"created_at"`
}
//---e-

//--- -s-
type Template struct {
	templates *template.Template
}
//
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
//---e-

//--- -s-
type Data struct {
	Comments []Comment
	Errors   []error
}
//---e-

func formatDateTime(d time.Time) string {
	if d.IsZero() {
		return ""
	}
	return d.Format("2006-01-02 15:04")
}

func main() {
    e := echo.New()

    // エラーハンドラの設定
    e.HTTPErrorHandler = func(err error, c echo.Context) {
        e.Logger.Error(err)
        c.Render(http.StatusInternalServerError, "error", err.Error()) //---コマンドメソッド エラーは標準出力
    }

    e.Renderer = &Template{
        templates: template.Must(template.New("").
            Funcs(template.FuncMap{
                "FormatDateTime": formatDateTime,
            }).ParseFS(templates, "templates/*")),
    }

    e.GET("/", func(c echo.Context) error {
        // ダミーのコメントデータ
        comments := []Comment{
            {Content: "This is a comment_test.", Created: time.Now()},
            {Content: "Another comment_test.", Created: time.Now().Add(-time.Hour)},
        }

        return c.Render(http.StatusOK, "index", Data{Comments: comments}) //---コマンドメソッド htmlにレンダリング
    })

	//---embedされたファイルシステムすべてを配信する設定(htmlとcss両方)-s-
	staticFs, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FileSystem(http.FS(staticFs))) //静的ファイルサーバを作成
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", fileServer))) //---ルーティングを設定

    e.Logger.Fatal(e.Start(":8989"))

}


