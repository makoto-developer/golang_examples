# Worker

- タスクを分散させて処理できる
  - 基本はラウンドロビン方式(workerの順番通りにメッセージが届く)
  - ACK(メッセージが届いたかどうかの確認)が保証されているので、万が一届かなかったら自動で再送される(キューに再配置される)
    - 他のworkerが見つかったらそちらに送りつける

```shell
# shell 1
go run worker.go
# shell 2
go run worker.go
# shell 3
# 何度も実行させてみるとshell 1とshell 2でそれぞれメッセージが届いていることがわかる
go run new_task.go
```