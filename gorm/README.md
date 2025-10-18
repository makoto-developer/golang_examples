# Gorm実装

GormとはGolang用のORM

作り途中で、説明用に作ったプロトタイプなのでツッコミはご遠慮ください。

# Requirement

# Use

- Docker (PostgreSQLサーバを起動するのに必要/ローカルに直接ホストしている人とかは不要)
- Gin(Web Framework)
- Gorm(ORMapper)

http serverをゼロから作るのは面倒だったのでGinを使っています。

# Setup

```shell
# Setupの手順に従ってPostgreSQLサーバを起動(※自前でホストしている人はスキップ)
git clone https://github.com/makoto-developer/docker-templates.git
cd docker-templates/postgresql-single-server/
docker compose up -d

# Goをインストール
mise i

# goのサーバを起動
go mod tidy
go run main.go
```

# References

公式サイト

- https://gorm.io/docs/

日本語はこちら

- https://gorm.io/ja_JP/docs/connecting_to_the_database.html

# コマンドサンプル集

注文を作る

```shell
curl -X POST http://localhost:8080/orders \
                                -H "Content-Type: application/json" \
                                -d '{
                              "order_item_group_id": 1,
                              "user_id": 100,
                              "amount": 15000,
                              "amount_without_tax": 13636,
                              "tax": 1364
                            }'

```

注文(複数の注文IDで)検索

```shell

# クエリなし（説明メッセージが返される）
curl "http://localhost:8080/orders"
# 複数のIDを指定して検索
curl "http://localhost:8080/orders?ids=1,2"
# レスポンスを整形して表示（jqコマンド使用）
curl "http://localhost:8080/orders?ids=1,2" | jq '.'
# 単一のIDで検索
curl "http://localhost:8080/orders?ids=1"
# verbose モード（詳細なレスポンス情報を表示）
curl -v "http://localhost:8080/orders?ids=1,2"
# レスポンスをファイルに保存
curl "http://localhost:8080/orders?ids=1,2" -o orders_response.json
```

注文(注文IDで)検索

```shell
# IDが1の注文を取得
curl "http://localhost:8080/orders/1"
# レスポンスを整形して表示（jqコマンド使用）
curl "http://localhost:8080/orders/1" | jq '.'
# IDが2の注文を取得
curl "http://localhost:8080/orders/2"
# IDが100の注文を取得
curl "http://localhost:8080/orders/100"
# verbose モード（詳細なレスポンス情報を表示）
curl -v "http://localhost:8080/orders/1"
# レスポンスをファイルに保存
curl "http://localhost:8080/orders/1" -o order_response.json
# HTTPステータスコードのみを表示
curl -w "%{http_code}\n" -o /dev/null -s "http://localhost:8080/orders/1"
# ヘッダー情報も含めて表示
curl -i "http://localhost:8080/orders/1"
```

ユーザIDで注文を検索

```shell
# ユーザーID 100 の全注文を取得
curl "http://localhost:8080/users/100/orders"
# レスポンスを整形して表示（jqコマンド使用）
curl "http://localhost:8080/users/100/orders" | jq '.'
# verbose モード（詳細なレスポンス情報を表示）
curl -v "http://localhost:8080/users/100/orders"
# レスポンスをファイルに保存
curl "http://localhost:8080/users/100/orders" -o user_orders_response.json
# HTTPステータスコードのみを表示
curl -w "%{http_code}\n" -o /dev/null -s "http://localhost:8080/users/100/orders"
# ヘッダー情報も含めて表示
curl -i "http://localhost:8080/users/100/orders"
```
