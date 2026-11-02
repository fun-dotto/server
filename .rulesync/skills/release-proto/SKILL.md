---
name: release-proto
description: release-proto ワークフローを手動実行し、app・web・admin-web の update-proto ワークフローを発火して proto の生成コードを更新する PR を作らせる。proto の変更をクライアントへ反映したいときに使う。
---

`.github/workflows/release-proto.yml` を `workflow_dispatch` で実行する。

## 手順

1. 実行するブランチを決める。指定がなければ `main` にする。クライアント側は常に fun-dotto/server の `main` から生成するため、通常は `main` で実行する。
2. ワークフローを実行する。

   ```bash
   gh workflow run release-proto.yml --ref main
   ```

3. 実行を特定して完了まで待つ。

   ```bash
   RUN_ID=$(gh run list --workflow release-proto.yml --limit 1 --json databaseId --jq '.[0].databaseId')
   gh run watch "$RUN_ID" --exit-status
   ```

4. 発火先の update-proto ワークフローの実行を確認する。

   ```bash
   for repo in app web admin-web; do
     gh run list --repo "fun-dotto/$repo" --workflow update-proto.yaml --limit 1
   done
   ```

## 結果の報告

- release-proto の実行 URL と結論（success / failure）を伝える。
- 各リポジトリの update-proto の状態と、作成・更新された PR（`chore/update-proto` ブランチ）があればその URL を伝える。
- 失敗した場合は `gh run view "$RUN_ID" --log-failed` でログを確認し、原因を伝える。トークンの発行に失敗した場合は、`APP_ID` / `APP_PRIVATE_KEY` の設定と、GitHub App が対象リポジトリにインストールされているかを確認するよう伝える。
