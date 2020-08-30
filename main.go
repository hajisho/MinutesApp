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

func setupRouter() *gin.Engine {

	r := gin.Default()

	dbInit() //データベースマイグレート

	// セッションの設定
	store := cookie.NewStore([]byte("secret"))
	// 静的ファイルのディレクトリを指定
	r.Static("dist", "./dist")
	// HTML ファイルのディレクトリを指定
	r.LoadHTMLGlob("./dist/public/*.html")

	r.Use(sessions.Sessions("mysession", store))
	// / に　GETリクエストが飛んできたらhandler関数を実行
	r.GET("/", returnMainPage)
	// /message に　GETリクエストが飛んできたらfetchMessage関数を実行
	r.GET("/message", fetchMessage)
	// ミーティング一覧を返す
	r.GET("/api_meetings", handleMeetings)
	// /add_messageへのPOSTリクエストは、handleAddMessage関数でハンドル
	r.POST("/add_message", handleAddMessage)
	// /update_messageへのPOSTリクエストは、handleUpdateMessage関数でハンドル
	r.POST("/update_message", handleUpdateMessage)
	// /update_messageへのPOSTリクエストは、handleDeleteMessage関数でハンドル
	r.POST("/delete_message", handleDeleteMessage)
	// ユーザー情報を返す
	r.GET("/user", fetchUserInfo)
	// ログインページを返す
	r.GET("/login", returnLoginPage)
	// ログイン動作を司る
	r.POST("/login", postLogin)
	//ユーザー登録ページを返す
	r.GET("/register", returnRegisterPage)
	//　ユーザー登録動作を司る
	r.POST("/register", postRegister)
	//セッション情報の削除
	r.GET("/logout", postLogout)
	// ミーティング一覧のページ
	r.GET("/meetings", returnMeetingsPage)

	r.GET("/entrance", returnEntrancePage)

	return r
}

func main() {

	//Temp="test"
	router := setupRouter()
	// サーバーを起動しています
	router.Run(":10000")
}

func returnMainPage(ctx *gin.Context) {
	//Cookieがなければログインページにリダイレクト　のつもり
	// 下記の関数、sessionCeckで確認しているから、実際には必要ないはず。要らなければ削除
	session := sessions.Default(ctx)
	user := session.Get("SessionID")
	if user == nil {
		ctx.Redirect(http.StatusSeeOther, "/entrance")
		ctx.Abort()
		return
	}
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "議事録", "header": "minuteHeader", "id": []string{"message"}})
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

// ResponseMeeting は、ミーティングがクライアントへ返される時の形式です。
type ResponseMeeting struct {
	Name string `json:"name"`
}

//ログインページのhtmlを返す
func returnLoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "Login and Register", "header": "loginHeader", "id": []string{"LoginAndRegister", "serverMessage"}})
}

//ユーザー登録ページのhtmlを返す
func returnRegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "Login and Register", "id": []string{"LoginAndRegister", "serverMessage"}})
}

func returnEntrancePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "Entrance", "header": "entranceHeader", "id": []string{"entrance", "serverMessage"}})
}

func returnMeetingsPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "Meetings", "header": "minuteHeader", "id": []string{"meetings"}})
}

//messagesに含まれるものを jsonで返す
func fetchMessage(ctx *gin.Context) {

	session := sessions.Default(ctx)

	if session.Get("SessionID") == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

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
	sessionID := session.Get("SessionID")

	if sessionID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if !(SessionExist(sessionID.(string))) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	user := getSessionUserID(ctx)

	//メッセージをデータベースへ追加
	dbInsert(req.Message, user.ID)

	ctx.JSON(http.StatusOK, gin.H{"success": true})

	return
}

// UpdateMessageRequest は、クライアントからのメッセージ追加要求のフォーマットです。
type UpdateMessageRequest struct {
	ID      string `json:"id"`
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

	session := sessions.Default(ctx)

	if session.Get("SessionID") == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	user := getSessionUserID(ctx)
	msg := dbGetOne(id)

	if user.ID != msg.UserID {
		// 権限がない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to privileges"})
		// 帰ることを忘れない
		return
	}

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

	session := sessions.Default(ctx)

	if session.Get("SessionID") == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	user := getSessionUserID(ctx)
	msg := dbGetOne(id)

	if user.ID != msg.UserID {
		// 権限がない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to privileges"})
		// 帰ることを忘れない
		return
	}

	//データベースにある指定されたメッセージを更新
	dbDelete(id)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// ユーザー自身の情報を返す
func fetchUserInfo(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("SessionID") == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	user := getSessionUserID(ctx)
	userInfo := ResponseUserPublic{
		ID:   user.ID,
		Name: user.Username,
	}
	ctx.JSON(http.StatusOK, userInfo)
}

//ログイン試行時にクライアントから送られてくるフォーマット
type userInfo struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

//登録動作
func postRegister(ctx *gin.Context) {
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
	return
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

	//セッション管理
	sessionID := createSession(user.Username)
	if sessionID == "0" {
		ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.Abort()
		return
	}

	//セッションにデータを格納する
	session := sessions.Default(ctx)
	session.Set("SessionID", sessionID)
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

//ログアウト処理
func postLogout(ctx *gin.Context) {

	//セッションからデータを破棄する
	session := sessions.Default(ctx)

	sessionID := session.Get("SessionID").(string)
	sessionDelete(sessionID)

	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/entrance")

}

func handleMeetings(ctx *gin.Context) {
	ms := getAllMeeting()
	ret := make([]ResponseMeeting, len(ms))
	for i, meeting := range ms {
		ret[i] = ResponseMeeting{
			Name: meeting.Name,
		}
	}
	ctx.JSON(http.StatusOK, ret)
}
