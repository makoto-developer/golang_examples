# RabbitMQ

## sample1

sendコマンドの引数の単語を送って、consumerがMQ経由で受信する簡単な実装

## sample2

Workerの実装。時間のかかる処理をタスクとしてMQで保存して、作業開始できる状態になったらconsumeする。

## sample3

Workerの実装。Ack(false)を設定すると処理中にworkerが落ちてもメッセージが再配置され再開可能になる。

