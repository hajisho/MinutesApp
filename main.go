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
)

func main() {

	router := gin.Default()
	// 静的ファイルのディレクトリを指定
	router.Static("dist", "./dist")
	// HTML ファイルのディレクトリを指定
	router.LoadHTMLGlob("./dist/public/*.html")

	dbInit() //データベースマイグレート

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
	messages := dbGetAll()
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

	//メッセージをデータベースへ追加
	dbInsert(req.Message)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
