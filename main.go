package main

//goはスペースが意味を持っているっぽい（調べてない
//pythonのインデントの感じ
//違うのはエラーなのにエラーと出力されないことがあること

import (
	//ginのインポート
	"net/http"

	"github.com/gin-gonic/gin"
	//"encoding/json"
	//"fmt"

	_ "github.com/mattn/go-sqlite3" //DBのパッケージだが、操作はGORMで行うため、importだけして使わない
)

// Message は、一つのメッセージの情報を持った構造体です。
type Message struct {
	//`json:"id"` を　`json: "id"`　にすると読み込めずにエラーなのだが、出力されないので気づかない
	//なので指定しているはずのに小文字にならないという勘違いが発生する
	// ID は、このメッセージの識別子として導入される値です。
	ID int `json:"id"`
	// Message は、このメッセージの文章です。
	Message string `json:"message"`
}

// Messages は、複数メッセージのスライスです。
type Messages []Message

// TODO データベースへ移行後削除
var messages Messages

//ID用
var count int

func addMessage(message string) {
	messages = append(messages, Message{count, message})
	count++
}

func main() {
	count = 0

	// TODO データベースの実装で削除
	messages = make([]Message, 0)

	//fmt.Printf("(%%#v) %#v\n", messages)

	router := gin.Default()
	// 静的ファイルのディレクトリを指定
	router.Static("dist", "./dist")
	// HTML ファイルのディレクトリを指定
	router.LoadHTMLGlob("./dist/public/*.html")
	// / に　GETリクエストが飛んできたらhandler関数を実行
	router.GET("/", handler)
	// /message に　GETリクエストが飛んできたらfetchMessage関数を実行
	router.GET("/message", fetchMessage)
	// /add_messageへのPOSTリクエストは、handleAddMessage関数でハンドル
	router.POST("/add_message", handleAddMessage)
	// サーバーを起動しています
	router.Run(":10000")
}

// 引数の型はデフォルトだと思います、引数名は任意でしょう
func handler(ctx *gin.Context) {
	// gin.H{}で、go ファイルの変数を HTML テンプレートに渡します。この例では何も渡していません。
	ctx.HTML(200, "index.html", gin.H{})
}

//messagesに含まれるものを jsonで返す
func fetchMessage(ctx *gin.Context) {
	ctx.JSON(200, messages)
}

// AddMessageRequest は、クライアントからのメッセージ追加要求のフォーマットです。
type AddMessageRequest struct {
	Message string `json:"message"`
}

func handleAddMessage(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(AddMessageRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.Message == "" {
		// メッセージがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'message' being empty"})
		// 帰ることを忘れない
		return
	}

	// TODO データベースへの蓄積
	addMessage(req.Message)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
