package main

//goはスペースが意味を持っているっぽい（調べてない
//pythonのインデントの感じ
//違うのはエラーなのにエラーと出力されないことがあること

import (
	//ginのインポート
	"github.com/gin-gonic/gin"
	//"encoding/json"
	//"fmt"
)

type Message struct {
	  //`json:"id"` を　`json: "id"`　にすると読み込めずにエラーなのだが、出力されないので気づかない
		//なので指定しているはずのに小文字にならないという勘違いが発生する
    Id   int    `json:"id"`
    Message string `json:"message"`
}

type Messages []Message
var messages Messages
//ID用
var count int

func addMessage(message string){
	messages = append(messages,Message{count,message})
	count++;
}


func main() {
	count=0
	addMessage("餃子")
	addMessage("チャーハン")

	//fmt.Printf("(%%#v) %#v\n", messages)

	router := gin.Default()
	// 静的ファイルのディレクトリを指定
	router.Static("dist", "./dist")
	// HTML ファイルのディレクトリを指定
	router.LoadHTMLGlob("./dist/public/*.html")
	// / に　GETリクエストが飛んできたらhandler関数を実行
	router.GET("/", handler)

	router.GET("/message", fetchMessage)
	// サーバーを起動しています
	router.Run(":10000")
}

// 引数の型はデフォルトだと思います、引数名は任意でしょう
func handler(ctx *gin.Context) {
	// gin.H{}で、go ファイルの変数を HTML テンプレートに渡します。この例では何も渡していません。
	ctx.HTML(200, "index.html", gin.H{})
}

//messagesに含まれるものを jsonで返す
func fetchMessage(ctx *gin.Context){
	ctx.JSON(200, messages)
}
