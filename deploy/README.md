# デプロイ手順

## 公開環境

- DigitalOcean Droplet
- Ubuntu 24.04 LTS
- Nginx：画面配信・Go APIへの転送
- systemd：Go APIの常時起動・自動起動
- 配置先：`/opt/space-now`
- Go API：`127.0.0.1:8000`

公開URL：

- デモ：http://143.198.80.71/web/
- APIドキュメント：http://143.198.80.71/docs/

現在の配置ブランチは `chore/deployment`。
通過一覧・詳細はモックを使用し、都市一覧とISS現在位置は実APIへ接続しています。
現在はHTTPで公開しています。

## 設定ファイル

- `nginx/space-now.conf`：Nginxの振り分け設定
- `systemd/space-now.service`：Go APIの起動設定

別のサーバーで利用する場合は、Nginxの `server_name` を変更してください。

## サーバーへの接続

手元のWSLで実行します。

```bash
ssh -i ~/.ssh/id_ed25519 root@143.198.80.71
```

SSHキーのパスは、各自の環境に合わせて変更してください。

## 初期配置

以下はサーバー側で実行します。
Goは、go.modで指定されたバージョン以上を公式配布からインストールしてください。
現在の環境では `/usr/local/go/bin/go` を使用しています。

```bash
apt update
apt install -y nginx git curl ca-certificates

git clone --branch chore/deployment --single-branch \
  https://github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017.git \
  /opt/space-now

cd /opt/space-now
/usr/local/go/bin/go build -o /opt/space-now/space-now .

adduser --system --group --no-create-home space-now

cp deploy/systemd/space-now.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now space-now

cp deploy/nginx/space-now.conf /etc/nginx/sites-available/space-now
ln -s /etc/nginx/sites-available/space-now /etc/nginx/sites-enabled/space-now
nginx -t
```

Nginxの設定確認が成功したら反映します。

```bash
systemctl enable --now nginx
systemctl reload nginx
```

## HTML・CSS・JS・ドキュメントの更新

サーバー側で実行します。

```bash
cd /opt/space-now
git pull --ff-only origin chore/deployment
```

これらのファイルだけの変更なら、Goの再ビルドやサービスの再起動は不要です。
ブラウザで再読み込みして確認してください。

## Goコードの更新

サーバー側で実行します。
まず新しい実行ファイルを作り、ビルド成功後に入れ替えます。

```bash
cd /opt/space-now
git pull --ff-only origin chore/deployment
/usr/local/go/bin/go build -o /opt/space-now/space-now.new .
```

ビルドが成功した場合だけ、次を実行します。

```bash
mv /opt/space-now/space-now.new /opt/space-now/space-now
systemctl restart space-now
systemctl status space-now --no-pager
```

## 動作確認

サーバー側で実行します。

```bash
systemctl is-active nginx
systemctl is-active space-now
curl -i http://127.0.0.1:8000/api/health
curl -i http://143.198.80.71/api/health
```

両サービスが `active`、healthが200と `{"status":"OK"}` を返すことを確認します。

ブラウザでも、地図・都市選択・ISS現在位置・ドキュメントへの往復を確認します。

## エラーの確認

```bash
journalctl -u space-now -n 50 --no-pager
tail -n 50 /var/log/nginx/error.log
```

## 今後の対応

- HTTPSの設定
- 詳細API完成後の通過一覧・詳細の実API接続
- mainへのマージ後、配置ブランチと本手順の更新