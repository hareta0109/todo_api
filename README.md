# TODO 管理 API

- 会社には管理会社と利用会社の２種類が、ユーザには管理者と一般ユーザの２種類が存在し、企業情報やユーザ情報の扱いに関する権限が分かれている。
- 管理会社のIDは1で固定である。
- タスクに関しては、編集者と閲覧者が存在し、編集者はタスクの追加・編集・閲覧が可能だが、閲覧者は閲覧のみ可能である。

## シードデータについて
### Company
- ID:1 管理会社
- ID:2 利用会社

### User
- ID:1 管理会社の管理者 password: super-admin
- ID:2 管理会社のユーザ password: super-user
- ID:3 利用会社の管理者 password: normal-admin
- ID:3 利用会社の一般編集者 password: normal-user
- ID:3 利用会社の一般閲覧者 password: normal-user

### シードデータの所在
`build/db/seeds/seeds.go` を参照

# requirements

- Docker 24.0.6
- go-task 3.30.1
- GoLang

# Usage

## Getting Started

```
# DBの初期化と起動(dockerコンテナが起動)
$ task db-init

# seedsデータをDBに登録する
$ task db-seeds

# internal/main.go を実行
# localhost:8080/swaggerにアクセスするとAPIのドキュメントが確認可能
$ task run  # serve at localhost:8080
```

## swagger の実行

```
# /docs にドキュメントが生成される
$ task swag
```

## linter の実行

```
# .golanci.yml を元に golangci-lint が実行される
$ task lint

# 修正まで行いたい場合
$ task lint-fix
```

## 開発に必要なツールのインストール

```
# go-task のインストール
$ go install github.com/go-task/task/v3/cmd/task@v3.30.1

# 必要なツールのインストール
$ task install-tools
```

## DB 関連のコマンド

```
# DBの初期化と起動(dockerコンテナが起動)
$ task db-init

# seedsデータをDBに登録する
$ task db-seeds

# DBのマイグレーション
$ task db-up

# DBのマイグレーションを一つ戻す
$ task db-down
```

## DB のマイグレーション

`/build/db/migrations` に sql ファイルを追加していく。

# 残タスク

- 環境変数への DB の接続情報の移し替え
- JWT 発行のシークレットキーの生成方法の検討
