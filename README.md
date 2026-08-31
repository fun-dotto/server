# Dotto Server

[![CI](https://github.com/fun-dotto/server/actions/workflows/ci.yml/badge.svg)](https://github.com/fun-dotto/server/actions/workflows/ci.yml)

## ローカル開発

### セットアップ

```bash
mise install   # go / atlas / gcloud / node / portless を導入
mise setup     # Go の依存関係を取得
```

### API サーバーの起動

各 API サーバーは mise のタスクから [portless](https://portless.sh/) 経由で起動する。
ポート番号の代わりに固定の `.localhost` URL (HTTPS / HTTP2) でアクセスできる。

| タスク | URL |
| --- | --- |
| `mise run academic-api` | https://academic-api.dotto.localhost |
| `mise run user-api` | https://user-api.dotto.localhost |
| `mise run announcement-api` | https://announcement-api.dotto.localhost |

portless は 4000-4999 の空きポートを環境変数 `PORT` で子プロセスに渡し、
そのポートへリバースプロキシする。サーバー側は `server.Addr()` が `PORT` を読むため、
`PORT` が未設定の場合 (直接 `go run` した場合や Cloud Run 以外) は従来どおり `:8080` で待ち受ける。

補助タスク:

| タスク | 内容 |
| --- | --- |
| `mise run portless:doctor` | プロキシ・DNS・CA 信頼状態の確認 |
| `mise run portless:list` | 登録済みルートの一覧 |
| `mise run portless:stop` | プロキシの停止 |

### 注意点

- 初回起動時に portless がローカル CA を生成して信頼ストアに登録し、443 番ポートを bind する。
  macOS / Linux では 443 を bind するために `sudo` へ自動昇格するのでパスワードを求められる。
  手動で信頼させたい場合は `portless trust` を実行する。
- Safari は `.localhost` を自動解決しないため、`portless hosts sync` で `/etc/hosts` に追記する。
- portless を使わず素の HTTP で起動したい場合は `go run ./cmd/academic-api` を直接実行する
  (`:8080` で待ち受ける)。`PORT=3000 go run ./cmd/academic-api` のように上書きも可能。
- ベース名は `portless.json` の `name` (= `dotto`) で定義している。git worktree を使うと
  ブランチ名がサブドメインの接頭辞として自動で付与され、worktree ごとに別 URL になる。

&copy; 2026 Dotto
