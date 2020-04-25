package nhk_api_test

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

const timeLayout = "15:04"

// Payload.Text用のテンプレート
var tpl = `
{{ range $key, $val := . -}}
日時： {{ $key }}
{{ range $index, $var := $val }}
番組： {{ $var.Title }}
出演者： {{ $var.Act }}
時間： {{ timeFmt $var.StartTime }} 〜 {{ timeFmt $var.EndTime }}
{{ end }}
------------------------------
{{ end }}
`

type payload struct {
	Channel   string `json:"channel"`
	UserName  string `json:"username,omitempty"`
	Text      string `json:"text"`
	IconURL   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

// Payload.Textの作成と設定
func (p *payload) setText(m map[string][]E1) {

	// テンプレート内で使用する関数の設定
	fMap := template.FuncMap{
		"timeFmt": timeFmt,
	}

	// テンプレートの割当
	tpl := template.Must(template.New("tpl").Funcs(fMap).Parse(tpl))

	// テンプレートからデータ出力
	var buffer bytes.Buffer
	if err := tpl.Execute(&buffer, m); err != nil {
		fmt.Println(err)
	}

	// 出力したデータを文字列化
	p.Text = buffer.String()
}

// テンプレートに渡す関数。
// 番組表APIの日時の形式を hh:mm にしてくれる。
func timeFmt(date time.Time) string {
	return date.Format(timeLayout)
}
