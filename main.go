package main

//goはスペースが意味を持っているっぽい（調べてない
//pythonのインデントの感じ
//違うのはエラーなのにエラーと出力されないことがあること

import (
	"strconv"
	//ginのインポート
	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func main() {

	//Temp="test"

	router := gin.Default()
	// 静的ファイルのディレクトリを指定
	router.Static("dist", "./dist")
	// HTML ファイルのディレクトリを指定
	router.LoadHTMLGlob("./dist/public/*.html")

	dbInit() //データベースマイグレート

	// セッションの設定
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// / に　GETリクエストが飛んできたらhandler関数を実行
	router.GET("/", returnMainPage)
	// /message に　GETリクエストが飛んできたらfetchMessage関数を実行
	router.GET("/message", fetchMessage)
	// /add_messageへのPOSTリクエストは、handleAddMessage関数でハンドル
	router.POST("/add_message", handleAddMessage)
	// /update_messageへのPOSTリクエストは、handleUpdateMessage関数でハンドル
	router.POST("/update_message", handleUpdateMessage)
	// /update_messageへのPOSTリクエストは、handleDeleteMessage関数でハンドル
	router.POST("/delete_message", handleDeleteMessage)
	// ログインページを返す
	router.GET("/login", returnLoginPage)
	// ログイン動作を司る
	router.POST("/login", postLogin)
	//ユーザー登録ページを返す
	router.GET("/register", returnRegisterPage)
	//　ユーザー登録動作を司る
	router.POST("/register", tempChallengeRegister)
	//セッション情報の削除
	router.GET("/logout", postLogout)

	router.GET("/entrance", returnEntrancePage)

	// サーバーを起動しています
	router.Run(":10000")
}

func returnMainPage(ctx *gin.Context) {
	//Cookieがなければログインページにリダイレクト　のつもり
	// 下記の関数、sessionCeckで確認しているから、実際には必要ないはず。要らなければ削除
	session := sessions.Default(ctx)
	user := session.Get("UserId")
	if user == nil {
		ctx.Redirect(http.StatusSeeOther, "/entrance")
		ctx.Abort()
		return
	}
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "議事録","header": "minuteHeader", "id": []string{"message"}})
}
// ResponseUserPublic は、公開ユーザー情報がクライアントへ返される時の形式です。
// JSON形式へマーシャルできます。
type ResponseUserPublic struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ResponseMessage は、メッセージがクライアントへ返される時の形式です。
// JSON形式へマーシャルできます。
type ResponseMessage struct {
	ID      uint               `json:"id"`
	AddedBy ResponseUserPublic `json:"addedBy"`
	Message string             `json:"message"`
}
//ログインページのhtmlを返す
func returnLoginPage(ctx *gin.Context){
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title":"Login and Register","header": "loginHeader","id":[]string{"LoginAndRegister","serverMessage"}})
}

//ユーザー登録ページのhtmlを返す
func returnRegisterPage(ctx *gin.Context){
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title":"Login and Register","id":[]string{"LoginAndRegister","serverMessage"}})
}

func returnEntrancePage(ctx *gin.Context){
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title":"Entrance","header": "entranceHeader","id":[]string{"entrance","serverMessage"}})
}

//messagesに含まれるものを jsonで返す
func fetchMessage(ctx *gin.Context) {
	messagesInDB := dbGetAll()
	// データベースに保存されているメッセージの形式から、クライアントへ返す形式に変換する
	messages := make([]ResponseMessage, len(messagesInDB))
	for i, msg := range messagesInDB {
		// TODO データベースでJOIN？
		user := getUserByID(msg.UserID)
		messages[i] = ResponseMessage{
			ID: msg.ID,
			AddedBy: ResponseUserPublic{
				ID:   msg.UserID,
				Name: user.Username,
			},
			Message: msg.Message,
		}
	}
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

	session := sessions.Default(ctx)
	user := getUser(session.Get("UserId").(string))

	//メッセージをデータベースへ追加
	dbInsert(req.Message, user.ID)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateMessageRequest は、クライアントからのメッセージ追加要求のフォーマットです。
type UpdateMessageRequest struct {
	ID string `json:"id"`
	Message string `json:"message"`
}

func handleUpdateMessage(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(UpdateMessageRequest)
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

	id, _ := strconv.Atoi(req.ID)
	//データベースにある指定されたメッセージを更新
	dbUpdate(id, req.Message)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteMessageRequest は、クライアントからのメッセージ追加要求のフォーマットです。
type DeleteMessageRequest struct {
	ID string `json:"id"`
}

func handleDeleteMessage(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(DeleteMessageRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.ID == "" {
		// IDがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'id' being empty"})
		// 帰ることを忘れない
		return
	}

	id, _ := strconv.Atoi(req.ID)
	//データベースにある指定されたメッセージを更新
	dbDelete(id)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

//ログイン試行時にクライアントから送られてくるフォーマット
type userInfo struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

//登録動作テストのための臨時関数
func tempChallengeRegister(ctx *gin.Context) {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'userId' or 'password' being empty"})
		// 帰ることを忘れない
		return
	}

	// DBにユーザーの情報を登録
	if err := createUser(req.UserId, req.Password); err != nil {
		// ログインIDがすでに使用されている
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "already use this id"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

//ログイン処理
func postLogin(ctx *gin.Context) {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'userId' or 'password' being empty"})
		// 帰ることを忘れない
		return
	}

	// 入力されたIDをもとにDBからレコードを取得
	user := getUser(req.UserId)

	if user.ID == 0 {
		// DBにユーザーの情報がない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not exist"})
		return
	}

	if err := comparePassword(user.Password, req.Password); err != nil {
		// パスワードが間違っている
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}

	//セッションにデータを格納する
	session := sessions.Default(ctx)
	session.Set("UserId", user.Username)
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

//ログアウト処理
func postLogout(ctx *gin.Context) {

	//セッションからデータを破棄する
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/entrance")

}
