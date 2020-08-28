package main

import (
	//"fmt"
	//"reflect"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

//正規のsession情報を格納する
var mainCookie string = " "
//別アカウントのテスト用
var subCookie string = " "
//ダミーのsession情報
//jsonStr := `{"UserId":"gadasgadsgadsggwrgjrdjbthgkmd","Password":"rhhrs65uhhenbeszrs4643"}`
var dummyCookie string = `mysession=MTU5ODA5MDg0NHxEdi1CQkFFQ180SUFBUkFCRUFBQVBmLUNBQUVHYzNSeWFXNW5EQWdBQmxWelpYSkpaQVp6ZEhKcGJtY01Id0FkWjJGa1lYTm5ZV1J6WjJGa2MyZG5kM0puYW5Ka2FtSjBhR2RyYldRPXzDVdeNdyqRk_UaOgI-QqjM_yvCiQA7swpbBWn7F7Ll6w==; Path=/; Expires=Mon, 21 Sep 2020 10:07:24 GMT; Max-Age=0`
//サーバーのルーティング
var router = setupRouter()

var GetMinutesPageRoute string = "/"
var EntranceRoute string = "/entrance"
var GetUserInfo string = "/user"
var GetMessageRoute string = "/message"
var PostMessageRoute string = "/add_message"
var UpdateMessageRoute string = "/update_message"
var DeleteMessageRoute string = "/delete_message"
var LoginRoute string = "/login"
var RegisterRoute string = "/register"
var LogoutRoute string = "/logout"

//エントランスページはセッション情報がなくても取得できる
func Test_entrancePage(t *testing.T) {
	//testRequestの結果を保存するやつ
	resp := httptest.NewRecorder()
  //テストのためのhttp request
	req, _ := http.NewRequest("GET", EntranceRoute, nil)
	//requestをサーバーに流して結果をrespに記録
	router.ServeHTTP(resp, req)

	//bodyを取り出し
	body, _ := ioutil.ReadAll(resp.Body)
	//ステータスコードは200のはず
	assert.Equal(t, 200, resp.Code)
	//titleはEntrance
	assert.Contains(t, string(body), "<title>Entrance</title>")
}

//loginページはセッション情報がなくても取得できる
func Test_loginPage(t *testing.T) {
	//testRequestの結果を保存するやつ
	resp := httptest.NewRecorder()
  //テストのためのhttp request
	req, _ := http.NewRequest("GET", LoginRoute, nil)
	//requestをサーバーに流して結果をrespに記録
	router.ServeHTTP(resp, req)

	//bodyを取り出し
	body, _ := ioutil.ReadAll(resp.Body)
	//ステータスコードは200のはず
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "<title>Login and Register</title>")
}

//エントランスページはセッション情報がなくても取得できる
func Test_registerPage(t *testing.T) {
	//testRequestの結果を保存するやつ
	resp := httptest.NewRecorder()
  //テストのためのhttp request
	req, _ := http.NewRequest("GET", RegisterRoute, nil)
	//requestをサーバーに流して結果をrespに記録
	router.ServeHTTP(resp, req)

	//bodyを取り出し
	body, _ := ioutil.ReadAll(resp.Body)
	//ステータスコードは200のはず
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "<title>Login and Register</title>")
}

//idとpasswordがそれぞれ８文字以上の英数字だと登録できる
func Test_canRegister_id_and_password_more8_and_alphanumeric(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test1234","Password":"qwer7890"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "success")

}

//userIdの重複不可
func Test_cntRegister_same_id_and_password(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test1234","Password":"qwer7890"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"already use this id`)

}

//パスワードが同じ場合は許す
func Test_canRegister_same_password(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test5678","Password":"qwer7890"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "success")

}

//提案
//７文字以下のuserIdは許さない
/*
func Test_cntRegister_id_less8(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test567","Password":"qwer7890"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"いい感じのエラー文`)

}
*/

//提案
//７文字以下のpasswordは許さない
/*
func Test_cntRegister_password_less8(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test5678","Password":"qwer789"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"いい感じのエラー文`)

}
*/

//提案
//英字のみのpasswordは許さない
/*
func Test_cntRegister_password_only_alphabet(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test5678","Password":"qwertest"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"いい感じのエラー文`)

}
*/

//提案
//数字のみのpasswordは許さない
/*
func Test_cntRegister_password_only_num(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test5678","Password":"11111111"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"いい感じのエラー文`)

}
*/

//登録済みのユーザーはログイン可能
func Test_canLogin_registered_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"UserId":"test1234","Password":"qwer7890"}`
	req, _ := http.NewRequest(
			"POST",
			LoginRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)
	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")
	mainCookie = resp.Header().Get("Set-Cookie")
}


//未登録のユーザーではログインできない
func Test_cntLogin_not_registered_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"UserId":"te344567","Password":"3wer3333"}`
	req, _ := http.NewRequest(
			"POST",
			LoginRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)
	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"user not exist`)
}


//ログインせずに議事録ページにはいけない
//リダイレクト
func Test_redirect_minutesPage_not_logined(t *testing.T){

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetMinutesPageRoute, nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, 303, resp.Code)
}


//必須
//登録されていないユーザー情報を持ったsessionではアクセスできない
/*
func Test_cntAccess_minutesPage_dummySession(t *testing.T){

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetMinutesPageRoute, nil)
	req.Header.Set("Cookie", dummyCookie)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "いい感じのエラー")
}
*/


//ログインしたなら議事録ページに行ける
func Test_canAccess_minutesPage_logined(t *testing.T){

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetMinutesPageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "<title>議事録</title>")
}


//ログアウト後にログインはできない
func Test_logout(t *testing.T){

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", LogoutRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)

	assert.Equal(t, 303, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	tempCookie := resp.Header().Get("Set-Cookie");

	req, _ = http.NewRequest("GET", GetMinutesPageRoute, nil)
	req.Header.Set("Cookie", tempCookie)

	router.ServeHTTP(resp, req)

	assert.Equal(t, 303, resp.Code)
}

//登録していないユーザーはメッセージを取得不可
func Test_cntGetMessge_not_logined_user(t *testing.T){

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", GetMessageRoute,nil)

	router.ServeHTTP(resp, req)

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"Bad Request`)

}


//登録済みのユーザーはメッセージを送信可能
func Test_canAddMessge_logined_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"message":"カシスオレンジ"}`
	req, _ := http.NewRequest(
			"POST",
			PostMessageRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	req, _ = http.NewRequest("GET", GetMessageRoute,nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body),`[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"カシスオレンジ"}]`)
}

//登録していないユーザーはメッセージを送信不可
func Test_cntAddMessge_not_logined_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"message":"ジン"}`
	req, _ := http.NewRequest(
			"POST",
			PostMessageRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"Bad Request`)

}

//テストのために他ユーザーセッションを取得
func Test_getSubCookie(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"subA2222","Password":"qwegds890"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	req, _ = http.NewRequest(
			"POST",
			LoginRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	subCookie = resp.Header().Get("Set-Cookie")

}

//異なるユーザーはメッセージを更新不可
func Test_cntUpdateMessge_different_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1","message":"ストロングゼロ"}`
	req, _ := http.NewRequest(
			"POST",
			UpdateMessageRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", subCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Malformed request due to privileges"`)

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute,nil)
	req.Header.Set("Cookie", subCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t,200, resp.Code)
	assert.Contains(t, string(body),`[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"カシスオレンジ"}]`)
	assert.NotContains(t, string(body),`[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"ストロングゼロ"}]`)
}


//登録していないユーザーはメッセージを更新不可
func Test_cntUpdateMessge_not_logined_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1","message":"ストロングゼロ"}`
	req, _ := http.NewRequest(
			"POST",
			UpdateMessageRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Bad Request"`)

}

//同じユーザーはメッセージを更新可能
func Test_canUpdateMessge_same_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1","message":"ストロングゼロ"}`
	req, _ := http.NewRequest(
			"POST",
			UpdateMessageRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute,nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t,200, resp.Code)
	assert.Contains(t, string(body),`[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"ストロングゼロ"}]`)
	assert.NotContains(t, string(body),`[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"カシスオレンジ"}]`)
}


//異なるユーザーはメッセージを削除不可
func Test_cntDeleteMessge_different_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1"}`
	req, _ := http.NewRequest(
			"POST",
			DeleteMessageRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", subCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Malformed request due to privileges"`)

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute,nil)
	req.Header.Set("Cookie", subCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t,200, resp.Code)
	assert.Contains(t, string(body),`[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"ストロングゼロ"}]`)
}


//登録していないユーザーはメッセージを削除不可
func Test_cntDeleteMessge_not_logined_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1"}`
	req, _ := http.NewRequest(
			"POST",
			DeleteMessageRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Bad Request"`)

}

//同じユーザーはメッセージを削除可能
func Test_canDeleteMessge_same_user(t *testing.T){

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1"}`
	req, _ := http.NewRequest(
			"POST",
			DeleteMessageRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute,nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t,200, resp.Code)
	assert.NotContains(t, string(body),`[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"ストロングゼロ"}]`)
}


//ログインしていないユーザーはユーザー情報は帰らない
func Test_cntGetUserInfo_not_logined_user(t *testing.T){

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetUserInfo, nil)

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Bad Request"`)

}

//ログインしているユーザーはユーザー情報が帰る
//セキュリティ的に大丈夫か？
func Test_canGetUserInfo_logined_user(t *testing.T){

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetUserInfo, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `{"id":1,"name":"test1234"}`)

}
