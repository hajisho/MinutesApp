//引用: http://vivacocoa.jp/go/gin/gin_firststep.php

package main

import (
	// Go の標準パッケージではなく、gin のパッケージをインポートします
	"github.com/gin-gonic/gin"
)

func main() {
	// gin の変数を定義しています
	router := gin.Default()
	// css などの静的ファイルのディレクトリを指定しています
	router.Static("dist", "./dist")
	// HTML ファイルのディレクトリを指定しています
	router.LoadHTMLGlob("./dist/public/*.html")
	// "/"ルートと handler 関数を関連づけています。handler という関数名は任意で付けた関数名です
	router.GET("/", handler)
	// サーバーを起動しています
	router.Run()
}

// 引数の型はデフォルトだと思います、引数名は任意でしょう
func handler(ctx *gin.Context) {
	// 200 の意味はわかりません。デフォルトではないかと思います
	// gin.H{}で、go ファイルの変数を HTML テンプレートに渡します。この例では何も渡していません。
	ctx.HTML(200, "index.html", gin.H{})
}
