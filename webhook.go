package nhk_api_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// 番組表APIのAPIキー
	apiKey = os.Getenv("API_KEY")
	// 探したい出演者(カンマ区切り)
	performer = os.Getenv("PERFORMER")
	// 探したい番組名(カンマ区切り)
	programName = os.Getenv("PROGRAM_NAME")
	// 探したい期間(Max7日)
	term = os.Getenv("GET_TERM")
	// WebhookのURL
	webhookURL = os.Getenv("WEBHOOK_URL")
	// Slackのチャンネル名
	channel = os.Getenv("CHANNEL")
)

const (
	// 日付のフォーマットProgram List APIのjson名で使用
	dateLayout = "2006-01-02"
	// URL
	programListURL = "https://api.nhk.or.jp/v2/pg/list/130/e1/%s.json?key=%s"
)

// Webhook メインロジック
func Webhook() {

	// 今日の日付取得
	// LambdaはUTCなので、9H補正
	now := time.Now()
	now = now.Add(time.Duration(9) * time.Hour)

	// 期間設定を数値に変換
	term, _ := strconv.Atoi(term)

	// 見たい番組表格納用マップ
	// K: 日付、V: E1構造体
	programMap := make(map[string][]E1)
	// API実行(3日先まで)
	for i := 0; i < term; i++ {
		// URIに渡す日付を計算
		d := now.AddDate(0, 0, i)
		date := d.Format(dateLayout)

		// 番組表APIの呼び出し
		reqURL := fmt.Sprintf(programListURL, date, apiKey)
		resp, err := http.Get(reqURL)
		if err != nil {
			log.Println("エラー発生")
			log.Println(fmt.Sprintf("Error:%s", err))
			log.Println(fmt.Sprintf("Request:%s", reqURL))
			return
		}

		if resp.StatusCode != 200 {
			log.Println("ステータスコード異常")
			log.Println(fmt.Sprintf("Status code:%v", resp.StatusCode))
			log.Println(fmt.Sprintf("Request:%s", reqURL))
			return
		}

		defer resp.Body.Close()

		// 番組表APIレスポンスのjson解析
		var p Program
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			log.Println("json解析でエラー発生")
			log.Println(err)
			return
		}

		// 見たい番組名・出演者があれば、pListに格納
		pList := make([]E1, 0)
		for _, e1 := range p.List.E1 {
			if isNoticeTarget(e1) {
				pList = append(pList, e1)
			}
		}

		// 見たい番組名・出演者がなければ、翌日へ。
		l := len(pList)
		if l == 0 {
			log.Printf("%vに、探してる番組・出演者はありませんでした。\n", date)
			continue
		}
		log.Printf("%vに、%v件の番組が見つかりました。\n", date, l)

		// マップに格納
		programMap[date] = pList
	}

	// payloadの生成
	payload := payload{
		Channel:   channel,
		UserName:  "nhk番組取得お知らせ",
		IconEmoji: ":ghost:",
		Text:      "探している番組は見つかりませんでした。",
	}

	// 見たい番組が見つかれば、メッセージ(Text)を更新。
	if len(programMap) > 0 {
		payload.setText(programMap)
	}

	// jsonのエンコード
	jsonByte, err := json.Marshal(payload)
	if err != nil {
		log.Println("jsonのエンコードでエラー発生")
		log.Println(err)
		return
	}

	// Webhookの呼び出し
	_, err = http.PostForm(
		webhookURL,
		url.Values{"payload": {string(jsonByte)}},
	)
	if err != nil {
		log.Println("Webhookの実行でエラー発生")
		log.Println(err)
		return
	}
}

// Webhookで通知する対象か判定
func isNoticeTarget(e1 E1) bool {

	// 出演者のチェック
	performerList := strings.Split(performer, ",")
	for _, v := range performerList {
		if strings.Contains(e1.Act, v) {
			return true
		}
	}

	// 番組のチェック
	programNameList := strings.Split(programName, ",")
	for _, v := range programNameList {
		if strings.Contains(e1.Title, v) {
			return true
		}
	}
	return false
}
