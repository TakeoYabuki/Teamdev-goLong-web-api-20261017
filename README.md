# Space now

ISSの現在位置と、選択した都市からの観測チャンスを確認できるWeb APIとデモアプリです。

RecursionのGo Web APIチーム開発課題として、3人で開発しました。ISSの通過予測に雲量予報を組み合わせ、観測時刻・方角・空の高さ・見やすさの目安をまとめて確認できます。

## 公開URL

- [デモアプリ](https://143.198.80.71/web/)
- [APIドキュメント（Swagger UI）](https://143.198.80.71/docs/)

## 主な機能

- ISSの現在位置を2D地図に表示し、約5秒間隔で更新
- ページを開いてから取得した直近40分の移動履歴を表示
- 観測都市の選択と、観測地点・世界全体への地図移動
- 次回の可視通過候補と開始までのカウントダウン
- 今後72時間の通過予測一覧
- 雲量・最大仰角・明るさによる見やすさの判定
- 見やすさ・最大仰角による絞り込み
- 開始時刻・最大仰角・雲量による並べ替え
- ページ切り替えと、見やすさ別の件数集計
- 通過1件の詳細と判定理由の表示
- 取得中・取得失敗・候補なしの状態表示

都市検索APIは固定47都市に対応し、デモ画面ではそのうち12都市を選択できます。

## 使い方

1. デモアプリで観測する都市を選択します。
2. 左のカードで、次回候補の開始時刻・方角・最大仰角・雲量を確認します。
3. 一覧で条件を絞り込み、並べ替えます。
4. 「詳細」を押すと、通過の開始・最大仰角・終了時点の情報や判定理由を確認できます。

次回候補は、見やすさが最も高い通過ではなく、対象期間で開始時刻が最も早い通過です。一覧の絞り込み条件は、次回候補と見やすさ別集計には適用されません。

## 使用技術

| 分野 | 技術 |
| --- | --- |
| バックエンド | Go、標準ライブラリの `net/http` |
| フロントエンド | HTML、CSS、JavaScript（フレームワーク不使用） |
| 地図 | Leaflet、NASA GIBSの夜間画像 |
| APIドキュメント | OpenAPI 3.0.3、Swagger UI |
| 公開環境 | DigitalOcean Droplet、Ubuntu 24.04 LTS |
| 配信・運用 | Nginx、systemd、Let's Encrypt、Certbot |

Nginxが画面ファイルを配信し、`/api/` へのリクエストを同じサーバー内のGo APIへ転送します。Go APIはsystemdで常時起動し、HTTPS証明書はCertbotで自動更新します。

## 利用している外部API・データ

| 提供元 | 用途 |
| --- | --- |
| [Where the ISS at?](https://wheretheiss.at/w/developer) | ISSの現在位置・高度・速度 |
| [N2YO](https://www.n2yo.com/api/) | ISSの可視通過予測 |
| [Open-Meteo](https://open-meteo.com/) | 都市の代表地点の雲量予報 |
| [NASA GIBS](https://nasa-gibs.github.io/gibs-api-docs/) | 地図の夜間画像 |

地図にはNASA GIBS / Suomi NPP VIIRSの `VIIRS_CityLights_2012` を使用しています。背景は2012年の夜間光画像で、現在の昼夜や天気を示すものではありません。

雲量データはOpen-Meteoから取得し、通過時刻への照合と見やすさの判定を本アプリで行っています。

## API一覧

| メソッド | パス | 内容 |
| --- | --- | --- |
| GET | `/api/health` | Go APIの起動確認 |
| GET | `/api/iss/now` | ISSの現在位置 |
| GET | `/api/locations/search` | 固定都市の検索・一覧 |
| GET | `/api/iss/passes` | 都市別の通過予測・次回候補・集計 |
| GET | `/api/iss/passes/{id}` | 保存された通過1件の詳細 |

パラメータ・応答項目・エラーの詳細は、[公開APIドキュメント](https://143.198.80.71/docs/) または [docs/openapi.yaml](docs/openapi.yaml) を参照してください。

### リクエスト例

```bash
# ISSの現在位置
curl "https://143.198.80.71/api/iss/now"

# 都市名の部分一致検索
curl --get "https://143.198.80.71/api/locations/search" \
  --data-urlencode "q=鹿児島"

# 東京の通過予測を開始時刻順で3件取得
curl "https://143.198.80.71/api/iss/passes?city=tokyo&page=1&limit=3"

# 最大仰角20度以上の通過を、最大仰角が高い順に取得
curl "https://143.198.80.71/api/iss/passes?city=tokyo&minElevation=20&sort=-maxElevation"
```

詳細を取得する場合は、先に一覧を取得し、応答の `passes[].id` または `nextCandidate.id` を使用します。次の `{id}` を実際に返されたIDへ置き換えてください。

```text
GET /api/iss/passes/{id}?city=tokyo
```

## 見やすさの判定

次の順序で判定します。

| 順序 | 条件 | 判定 |
| --- | --- | --- |
| 1 | 雲量70%以上、または最大仰角20度未満 | 条件が厳しい（poor） |
| 2 | 1に該当せず、雲量または明るさが不明 | 判定不明（unknown） |
| 3 | 雲量30%未満・最大仰角40度以上・明るさ−2等級以下をすべて満たす | 良好（good） |
| 4 | それ以外 | 通常（fair） |

雲量は、最大仰角となる時刻に最も近い予報を使います。時刻差が30分を超える場合や値が欠けている場合は、不明として扱います。見やすさは観測条件の目安で、実際の視認を保証するものではありません。

## キャッシュと詳細データ

| データ | 保存期間・更新方法 |
| --- | --- |
| ISS現在位置 | 5秒間キャッシュ |
| 通過予測 | 都市ごとに1時間キャッシュ |
| 雲量予報 | 都市ごとに1時間キャッシュ |
| 通過詳細 | 一覧の指定ページ・次回候補を保存。同じ都市・通過IDも上書きし、最後の保存・更新から24時間保持 |

キャッシュ期限後のリクエストで、外部APIから再取得します。キャッシュと詳細はメモリ上に保存するため、Go APIの再起動で失われます。

古い一覧を表示したまま、別のリクエストで保存内容が更新された場合は、画面の一覧と詳細の雲量・判定が異なることがあります。その場合は一覧を再取得してください。

## ローカルでの起動

以下は、WSL・Linux・macOSのbash / zshでの手順です。

### 必要なもの

- Git
- Go 1.27.2以上（[go.mod](go.mod) の指定を満たすバージョン）
- N2YOのAPIキー（通過予測を試す場合）
- ブラウザー

### 1. リポジトリを取得

```bash
git clone https://github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017.git
cd Teamdev-goLong-web-api-20261017
```

### 2. 環境変数を設定

```bash
cp .env.example .env
```

`.env` を編集し、自分のN2YO APIキーを設定します。

```dotenv
PORT=8000
N2YO_API_KEY=取得したAPIキー
```

`.env` はGitの管理対象外です。APIキーをREADMEやJavaScriptへ記載しないでください。

### 3. Go APIを起動

```bash
set -a
source .env
set +a
go run .
```

Goは `.env` を自動で読み込みません。上の操作で環境変数として読み込んでから起動します。

APIは `127.0.0.1:8000` で待ち受けます。Goサーバー自体は画面用の静的ファイルを配信しません。

### 4. APIを確認

別のターミナルで実行します。

```bash
curl -i "http://localhost:8000/api/health"
curl -i "http://localhost:8000/api/iss/now"
curl -i "http://localhost:8000/api/iss/passes?city=tokyo&page=1&limit=3"
```

healthが `200 OK` と `{"status":"OK"}` を返せばGo APIは応答しています。N2YOのキーが未設定の場合、通過一覧は502になります。health・都市検索・ISS現在位置はN2YOキーなしでも利用できます。

### 5. デモ画面をローカルで確認

公開用の [web/app.js](web/app.js) は、同じホストのAPIへ接続する設定です。

```javascript
const API_BASE_URL = "";
```

VSCodeのLive Serverなどでフロントを別ポートから開く場合は、ローカル確認中のみ次に変更します。

```javascript
const API_BASE_URL = "http://localhost:8000";
```

Go APIを起動したまま、リポジトリをVSCodeで開き、`web/index.html` をLive Serverで表示してください。`file://` で直接開く手順ではありません。

APIの `PORT` を変更した場合は、上の接続先ポートも揃えてください。公開用の変更をコミットする前に、`API_BASE_URL` を空文字へ戻します。

ローカルと公開環境の接続先を手動で書き換えずに切り替える仕組みは、今後の改善対象です。

### 環境変数一覧

| 名前 | 用途 | 未設定時 |
| --- | --- | --- |
| `PORT` | Go APIの待受ポート | 8000 |
| `N2YO_API_KEY` | N2YOの通過予測取得 | 通過予測を取得できない |

### ドキュメントのローカル確認

リポジトリのルートをLive Serverで配信し、`/docs/index.html` を開きます。

OpenAPIの `servers.url` は `/` のため、「Try it out」の送信先はドキュメントを開いたホストです。Live ServerだけではGo APIへ転送されません。APIを試す場合はローカルのGo APIをcurlで呼ぶか、公開ドキュメントを使用してください。Nginxで画面とAPIを同じホストにまとめる構成については、[デプロイ手順](deploy/README.md) を参照してください。

## ビルド・テスト

リポジトリのルートで実行します。APIを手動で起動する必要はありません。

```bash
go build ./...
go vet ./...
go test ./...
go test -race ./...
```

## ディレクトリ構成

| パス | 内容 |
| --- | --- |
| `main.go` | APIルーティングとサーバー起動 |
| `internal/client/` | 外部APIへの通信 |
| `internal/handler/` | HTTPリクエストの検証・応答 |
| `internal/service/` | 判定・時刻照合・絞り込み・並べ替え・キャッシュ |
| `internal/model/` | データ構造・固定都市マスター |
| `web/` | デモアプリのHTML・CSS・JavaScript |
| `docs/` | OpenAPI定義とSwagger UI |
| `deploy/` | デプロイ手順、Nginx・systemd・証明書更新の設定 |

## 主な担当

担当は主担当の区分です。仕様の合意、各担当部分のテスト、PRレビュー、結合確認はチームで行いました。

| メンバー | 主な担当 |
| --- | --- |
| Takeo（A） | ISS現在位置APIの開発、外部APIとの通信、5秒キャッシュ、N2YOによる通過予測の取得・データ変換、関連テスト、PRレビュー |
| hiroki（B） | 雲量予報、時刻照合、見やすさ判定、共通JSON・エラー・CORS、通過一覧の絞り込み・並べ替え・ページ分け・集計、キャッシュ・詳細保存 |
| レイナ（C） | 固定都市検索API、Leaflet地図、都市選択、カウントダウン、一覧操作・詳細展開・集計・通信状態の表示、APIドキュメント画面、デプロイ |

## デプロイ

公開環境の初期設定、更新、証明書の自動更新、ログの確認方法は、[deploy/README.md](deploy/README.md) に記載しています。

本番のGo APIには、systemd経由で `N2YO_API_KEY` を渡す設定も必要です。リポジトリ内の `deploy/systemd/space-now.service` には現在その読み込み設定が含まれていません。別サーバーへ配置する場合は、秘密情報をGitへ登録せず、環境変数ファイルと `EnvironmentFile=` などを追加してください。

## 今後の改善・機能追加

### リファクタリング

- JavaScriptのAPI通信・画面更新・地図操作の役割を整理する
- ローカルと公開環境のAPI接続先の切り替えを簡単にする
- 通過開始時刻に同じ候補が返った場合の再取得処理を改善する

### 機能追加

- 地図クリックによる任意の観測地点の指定
- 小惑星の接近情報の表示
- ロケットの打ち上げ予定の表示

地図クリックによる地点指定は別ブランチで試作しており、公開中のmainには含まれていません。小惑星と打ち上げ予定は、当初の企画で検討した今後の追加機能です。
