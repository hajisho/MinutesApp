package main

//goはスペースが意味を持っているっぽい（調べてない
//pythonのインデントの感じ
//違うのはエラーなのにエラーと出力されないことがあること

import (
	//ginのインポート
	"net/http"
	"github.com/gin-gonic/gin"

)

//フロントエンドのログイン動作をテストするために作った臨時のグローバル変数
//バックエンドが完成しだい消してください
var Temp string

func main() {

	//Temp="test"

	router := gin.Default()
	// 静的ファイルのディレクトリを指定
	router.Static("dist", "./dist")
	// HTML ファイルのディレクトリを指定
	router.LoadHTMLGlob("./dist/public/*.html")

	dbInit() //データベースマイグレート

	// / に　GETリクエストが飛んできたらhandler関数を実行
	router.GET("/", mainPage)
	// /message に　GETリクエストが飛んできたらfetchMessage関数を実行
	router.GET("/message", fetchMessage)
	// /add_messageへのPOSTリクエストは、handleAddMessage関数でハンドル
	router.POST("/add_message", handleAddMessage)
  // ログインページを返す
	router.GET("/login",returnLoginPage)
  // ログイン動作を司る
	router.POST("/login",tempChallengeLogin)
  //セッション情報の削除のつもり
	router.GET("/logout",tempDeleteCookie)

	// サーバーを起動しています
	router.Run(":10000")
}


func mainPage(ctx *gin.Context) {
	//Cookieがなければログインページにリダイレクト　のつもり
	if(Temp==""){
		ctx.Redirect(http.StatusMovedPermanently, "/login")
  	ctx.Abort()
  }
	ctx.HTML(http.StatusOK, "index.html", gin.H{"title":"議事録","id":[]string{"message"}})
}

//messagesに含まれるものを jsonで返す
func fetchMessage(ctx *gin.Context) {
	messages := dbGetAll()
	ctx.JSON(http.StatusOK, messages)
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

//ログインページのhtmlを返す
func returnLoginPage(ctx *gin.Context){
	ctx.HTML(http.StatusOK, "index.html", gin.H{"title":"loginページ","id":[]string{"login"}})
}

//ログイン試行時にクライアントから送られてくるフォーマット
type userInfo struct {
	UserId string `json:"userId"`
	Password string `json:"password"`
}

//ログイン動作を司る
//クライアント動作確認のための仮関数
func tempChallengeLogin(ctx *gin.Context){
	// POST bodyからメッセージを獲得
	req := new(userInfo)
	err := ctx.BindJSON(req)

	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.UserId == "" || req.Password == "" {
		// メッセージがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'message' being empty"})
		// 帰ることを忘れない
		return
	}

	//Cookieセットのイメージ
	//本来はクライアント側にcookieが帰る
	Temp = req.UserId
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

//ログアウト動作のつもり
//動作未確認
func tempDeleteCookie(ctx *gin.Context){
	Temp = ""
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
