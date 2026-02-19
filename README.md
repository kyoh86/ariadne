# minecraft-client

Fabric + Modrinth 前提で、Minecraft プロファイル作成と Mod 更新を楽にするための CLI (`ariadne`) です。

## できること

- プロファイル作成 (`profile create`)
- プロファイル一覧 (`profile list`)
- アクティブプロファイル設定 (`profile use`)
- プロファイル削除 (`profile drop`)
- Modrinth の Mod 追加 (`mod add`)
- Mod 一覧 (`mod list`)
- Mod 削除 (`mod drop`)
- Modrinth の Mod 検索 (`mod search`)
- ShaderPack 追加 (`shader add`)
- ShaderPack 一覧 (`shader list`)
- ShaderPack 削除 (`shader drop`)
- Minecraft バージョンに合う Mod の同期 (`sync`)
- Minecraft バージョン更新と再同期 (`upgrade-mc --sync`)
- `profile create` 時に Fabric Loader を launcher-dir 配下へ初期化し、`launcher_profiles.json` を更新

## ビルド

```bash
go build -o bin/ariadne ./cmd/ariadne
```

## 使い方

1. デフォルトの game-dir ルートを設定（初回のみ）

```bash
./bin/ariadne config set-game-dir-root /mnt/d/minecraft
./bin/ariadne config set-launcher-dir /mnt/d/minecraft
```

2. プロファイル作成

```bash
./bin/ariadne profile create 1.21.5
```

プロフィール名は `--name` で指定できます。省略時は `fabricmc-<mc-version>` が使われます。
例: `./bin/ariadne profile create 1.21.5 --name main`
`--game-dir` を省略すると `<game-dir-root>/<name>` が使われます。
個別に上書きしたい場合だけ `--game-dir` を指定してください。
`profile create` は Fabric installer を実行するため `java` が必要です。

3. Mod を追加（slug または project ID）

```bash
./bin/ariadne profile use main
./bin/ariadne mod add fabric-api
./bin/ariadne mod add sodium
./bin/ariadne mod add lithium
./bin/ariadne mod list
./bin/ariadne mod drop sodium
# active profileを使わないなら:
./bin/ariadne mod add fabric-api --profile main
./bin/ariadne mod list --profile main
./bin/ariadne mod drop sodium --profile main
```

ShaderPack を管理する場合:

```bash
./bin/ariadne shader add bsl-shaders
./bin/ariadne shader list
./bin/ariadne shader drop bsl-shaders
# active profileを使わないなら:
./bin/ariadne shader add bsl-shaders --profile main
```

必要なModを探す場合:

```bash
./bin/ariadne mod search minimap
```

4. Mod を同期

```bash
./bin/ariadne sync
# active profileを使わないなら:
./bin/ariadne sync --profile main
```

5. Minecraft バージョン更新（更新直後の運用向け）

```bash
./bin/ariadne upgrade-mc 1.21.6 --sync
# active profileを使わないなら:
./bin/ariadne upgrade-mc 1.21.6 --profile main --sync
```

6. プロファイル一覧と削除

```bash
./bin/ariadne profile list
./bin/ariadne profile use main
./bin/ariadne profile drop main
# game-dir も消す場合
./bin/ariadne profile drop main --delete-game-dir
```

## 保存されるデータ (XDG)

- `$XDG_CONFIG_HOME/ariadne/config.json` (`XDG_CONFIG_HOME` 未設定時は `~/.config/ariadne/config.json`)
  - 既定の game-dir ルート / launcher-dir 設定
- `$XDG_DATA_HOME/ariadne/profiles/<name>.json` (`XDG_DATA_HOME` 未設定時は `~/.local/share/ariadne/profiles/<name>.json`)
  - プロファイル定義（MC バージョン、Fabric Loader、Mod 一覧）
- `$XDG_DATA_HOME/ariadne/locks/<name>.json` (`XDG_DATA_HOME` 未設定時は `~/.local/share/ariadne/locks/<name>.json`)
  - 最終同期で使った Mod ファイル情報
- `$XDG_DATA_HOME/ariadne/cache/fabric-installer-<version>.jar` (`XDG_DATA_HOME` 未設定時は `~/.local/share/ariadne/cache/fabric-installer-<version>.jar`)
  - Fabric installer のキャッシュ

`sync` 実行時は lock 情報を使って、古くなった管理対象 jar を削除します。
