# JWT認証実装

JWT (JSON Web Token) 認証の実装サンプル。

## 技術スタック

- **Go**: 1.25+
- **Gin**: Webフレームワーク
- **golang-jwt/jwt**: JWT生成・検証ライブラリ
- **bcrypt**: パスワードハッシュ化

## 機能

- ユーザー登録 (パスワードハッシュ化)
- ログイン (JWT発行)
- JWT検証ミドルウェア
- リフレッシュトークン
- 保護されたエンドポイント

## ディレクトリ構成

```
jwt-auth/
├── server/
│   ├── handler/        # HTTPハンドラ
│   │   ├── auth.go     # 認証ハンドラ (登録、ログイン)
│   │   └── user.go     # ユーザーハンドラ (保護されたエンドポイント)
│   ├── middleware/     # ミドルウェア
│   │   └── auth.go     # JWT検証ミドルウェア
│   ├── model/          # データモデル
│   │   ├── user.go     # ユーザーモデル
│   │   └── token.go    # トークンモデル
│   └── util/           # ユーティリティ
│       └── jwt.go      # JWT生成・検証
├── main.go             # エントリーポイント
├── .env.example        # 環境変数サンプル
├── .gitignore
└── mise.toml
```

## セットアップ

### 1. 依存関係のインストール

```bash
mise install
go mod tidy
```

### 2. 環境変数の設定

```bash
cp .env.example .env
# .envを編集してJWT_SECRETを設定
```

### 3. サーバー起動

```bash
go run main.go
もしくは
mise run server
```

サーバーは `http://localhost:22500` で起動します。

## 使い方

### ユーザーを登録

```bash
curl -X POST http://localhost:22500/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### login

```bash
curl -X POST http://localhost:22500/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

レスポンス:

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600
}
```

### 保護されたエンドポイントへのアクセス

```bash
curl http://localhost:22500/api/users/me \
  -H "Authorization: Bearer <access_token>"
```

### リフレッシュトークン

```bash
curl -X POST http://localhost:22500/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```

## セキュリティの設計

- パスワードはbcryptでハッシュ化
- JWTシークレットは環境変数で管理
- アクセストークンの有効期限: 1時間
- リフレッシュトークンの有効期限: 7日間
- HTTPS通信を推奨 (本番環境)

## 環境変数

| 変数名              | 説明                             | デフォルト |
| ------------------- | -------------------------------- | ---------- |
| PORT                | サーバーポート                   | 22500      |
| JWT_SECRET          | JWT署名用シークレット            | (必須)     |
| JWT_ACCESS_EXPIRES  | アクセストークン有効期限(秒)     | 3600       |
| JWT_REFRESH_EXPIRES | リフレッシュトークン有効期限(秒) | 604800     |

## Test

```bash
mise run test
```

## References

- [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- [Gin Web Framework](https://gin-gonic.com/)
- [JWT.io](https://jwt.io/)
