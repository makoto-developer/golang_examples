user@usernoMacBook-Pro ~/work/repositories/github.com/makoto-developer/golang_examples/gorm (main)> go mod init github.com/makoto-developer/golang_examples/gorm
go: creating new go.mod: module github.com/makoto-developer/golang_examples/gorm
user@usernoMacBook-Pro ~/work/repositories/github.com/makoto-developer/golang_examples/gorm (main)> vi main.go

```go
package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func main() {
  // loggerとrecoveryミドルウェア付きGinルーター作成
  r := gin.Default()

  // 簡単なGETエンドポイント定義
  r.GET("/ping", func(c *gin.Context) {
    // JSONレスポンスを返す
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  // ポート8080でサーバー起動（デフォルト）
  // 0.0.0.0:8080（Windowsではlocalhost:8080）で待機
  r.Run()
}

```
user@usernoMacBook-Pro ~/work/repositories/github.com/makoto-developer/golang_examples/gorm (main)> go mod tidy
user@usernoMacBook-Pro ~/work/repositories/github.com/makoto-developer/golang_examples/gorm (main)> go get -u gorm.io/gorm
