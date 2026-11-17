---
name: create-pr
description: PR テンプレートに従って、日本語のタイトルと本文で Draft PR を作成し、自分自身をアサインする。PR の作成を頼まれたときに使う。
---

現在のブランチから、PR テンプレートに従った Draft PR を作成する。

## 手順

1. 変更内容を把握する。

   ```bash
   git fetch origin main
   git log --oneline origin/main..HEAD
   git diff origin/main...HEAD
   ```

   タイトルと本文は、コミットメッセージではなく、ベースブランチとの差分全体を元に書く。コミットメッセージは途中経過を含み、最終的な変更内容と一致しないことがあるため、参考に留める。

   未コミットの変更があれば、先に `commit` スキルでコミットする。
2. ブランチをプッシュする。

   ```bash
   git push -u origin HEAD
   ```

3. PR テンプレートを読む。`.github/pull_request_template.md`、`.github/PULL_REQUEST_TEMPLATE.md`、`.github/PULL_REQUEST_TEMPLATE/` の順に探す。
4. タイトルと本文を日本語で書く。
   - 差分全体（`git diff origin/main...HEAD`）を読んで、最終的な変更内容を書く
   - タイトル: PR 全体の変更内容を 1 文で簡潔に書く
   - 本文: テンプレートの構成に従って書く
     - テンプレートにあるセクションは、書く内容がなくても削除しない。見出しは残し、空欄のままにする
     - テンプレートのコメント（`<!-- -->`）やチェックボックスもそのまま残し、該当するものだけチェックする
     - テンプレートがない場合は「概要」「変更内容」のセクションで書く
5. Draft PR を作成し、自分自身をアサインする。

   ```bash
   gh pr create --draft --base main --assignee @me --title "<タイトル>" --body-file <本文ファイル>
   ```

   `gh` が使えない環境では GitHub の MCP ツールで Draft PR を作成し、ログイン中のユーザー（`get_me` で取得）を assignee に設定する。
6. 作成した PR の URL を報告する。
