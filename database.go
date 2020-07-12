package main

import(
  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3" //DBのパッケージだが、操作はGORMで行うため、importだけして使わない
)

/*
gorm.Modelの中身
カラム
  ・id
  ・created_at
  ・updated_at
  ・deleted_at
*/
/*外部からカラムを参照するときは
id → ID
created_at → CreatedAt
updated_at → UpdatedAt
deleted_at → DeletedAt
*/
// テーブル名：messages -->　テーブル名は自動で複数形になる
type Message struct{
  gorm.Model
  Message string
  MeetingID int
  UserID int
}

type User struct {
  gorm.Model
  Username string `gorm:"unique;not null"`
  Password string
}

type Meeting struct {
  gorm.Model
  Name string
}

type Entry struct {
  gorm.Model
  MeetingID int
  UserID int
}
/*
DBの内容
(ID,作成日,更新日,削除日のカラムは全てに入っている)
・ユーザー
  ・ユーザーネーム（ログインID）
  ・パスワード（暗号化したもの）
・会議
  ・会議名
・メッセージ
  ・内容
  ・会議ID
  ・ユーザーID
・エントリー
  ・会議ID
  ・ユーザーID
*/

//DBマイグレート
//main関数の最初でdbInit()を呼ぶことでデータベースマイグレート
func dbInit(){
  db, err := gorm.Open("sqlite3", "minutes.sqlite3") //第一引数：使用するDBのデバイス。第二引数：ファイル名
  if err != nil{
    panic("データベース開ません(dbinit)")
  }
  db.AutoMigrate(&User{}, &Message{}, Meeting{}, &Entry{}) //ファイルがなければ、生成を行う。すでにあればマイグレート。すでにあってマイグレートされていれば何も行わない
  defer db.Close()
}

//DB追加
//追加したいメッセージは、dbInsert(message.Message)のような感じで呼べば追加される
func dbInsert(message string){
  db, err := gorm.Open("sqlite3", "minutes.sqlite3")
  if err != nil{
    panic("データベース開ません(dbInsert)")
  }
  db.Create(&Message{Message: message})
  defer db.Close()
}

//DB全取得
//dbGetAll()と呼ぶことで、データベース内の全てのMessageオブジェクトが返される
func dbGetAll() []Message{
  db, err := gorm.Open("sqlite3", "minutes.sqlite3")
  if err != nil{
    panic("データベース開ません(dbGetAll)")
  }
  var messages []Message
  db.Order("created_at desc").Find(&messages) //db.Find(&messages)で構造体Messageに対するテーブルの要素全てを取得し、それをOrder("created_at desc")で新しいものが上に来るように並び替えている
  db.Close()
  return messages
}

//DB一つ取得
//idを与えることで、該当するMessageオブジェクトが一つ返される
func dbGetOne(id int) Message{
  db, err := gorm.Open("sqlite3", "minutes.sqlite3")
  if err != nil{
    panic("データベース開ません(dbGetOne)")
  }
  var message Message
  db.First(&message, id)
  db.Close()
  return message
}

//DB更新
//idとmessageを与えることで、該当するidのMessageオブジェクトのMessageが更新される
func dbUpdate(id int, update_message string){
  db, err := gorm.Open("sqlite3", "minutes.sqlite3")
  if err != nil{
    panic("データベース開ません(dgUpdate)")
  }
  var message Message
  db.First(&message, id)
  message.Message = update_message
  db.Save(&message)
  db.Close()
}

//DB削除
//指定したidのMessageオブジェクトが削除される
func dbDelete(id int){
  db, err := gorm.Open("sqlite3", "minutes.sqlite3")
  if err != nil{
    panic("データベース開ません(dbDelete)")
  }
  var message Message
  db.First(&message, id)
  db.Delete(&message)
  db.Close()
}

// ユーザー登録処理
func createUser(username string, password string) error {
  db, err := gorm.Open("sqlite3", "minutes.sqlite3")
  if err != nil{
    panic("データベース開ません(dgUpdate)")
  }
  defer db.Close()
  // Insert処理
  if err := db.Create(&User{Username: username, Password: password}).Error; err != nil {
    return err
  }
  return nil
}

// ユーザーネーム(ログインID)を指定してそのユーザーのレコードを取ってくる
// 使用時
// user := getuser(userId)
// ログインID: user.Username パスワード: user.Password
func getUser(username string) User {
  db, err := gorm.Open("sqlite3", "minutes.sqlite3")
  if err != nil{
    panic("データベース開ません(dgUpdate)")
  }
  defer db.Close()
  var user User
  db.Where("username = ?", username).First(&user)
  return user
}
