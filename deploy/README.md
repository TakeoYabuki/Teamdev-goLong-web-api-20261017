# デプロイ手順

## 公開環境

- DigitalOcean Droplet
- Ubuntu 24.04 LTS
- Nginx：画面配信・Go APIへの転送
- systemd：Go APIの常時起動・自動起動
- 配置先：`/opt/space-now`
- Go API：`127.0.0.1:8000`
- 配置ブランチ：`chore/deployment`

公開URL：

- デモ：https://143.198.80.71/web/
- APIドキュメント：https://143.198.80.71/docs/

通過一覧・詳細はモックを使用し、都市一覧とISS現在位置は実APIへ接続しています。

HTTPSにはLet’s EncryptのIPアドレス用証明書を使用しています。
Certbotで自動更新し、更新後はdeploy hookでNginxを再読み込みします。

## 設定ファイル

- `nginx/space-now.conf`：HTTPS・画面配信・API転送の設定
- `systemd/space-now.service`：Go APIの起動設定
- `letsencrypt/reload-nginx`：証明書更新後のNginx再読み込み

別のサーバーで利用する場合は、本手順とNginx設定内のIPアドレス、
証明書のパスを変更してください。

証明書・秘密鍵はGitへ登録しません。

## サーバーへの接続

手元のWSLで実行します。

```bash
ssh -i ~/.ssh/id_ed25519 root@143.198.80.71
```

SSHキーのパスは、各自の環境に合わせて変更してください。

## 初期配置

以下は、新しいサーバーへ初めて配置するときの手順です。
すでに稼働しているサーバーでは繰り返し実行しません。

以降のコマンドは、特に記載がなければサーバー側で実行します。

### 1. 必要なソフトとGoを準備

```bash
apt update
apt install -y nginx git curl ca-certificates
```

Goは、go.modで指定されたバージョンを満たす公式配布版をインストールします。
現在の環境はGo 1.27.3で、`/usr/local/go/bin/go` を使用しています。

```bash
/usr/local/go/bin/go version
```

### 2. コードの取得とビルド

```bash
git clone --branch chore/deployment --single-branch \
  https://github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017.git \
  /opt/space-now

cd /opt/space-now
/usr/local/go/bin/go build -o /opt/space-now/space-now .
```

ビルドが成功したことを確認してから、次へ進みます。

### 3. Go APIをサービスとして登録

```bash
adduser --system --group --no-create-home space-now

cp deploy/systemd/space-now.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now space-now
systemctl status space-now --no-pager
```

`active (running)` と表示されることを確認します。

### 4. 証明書取得用のHTTP設定

リポジトリのNginx設定はHTTPS用です。
証明書がまだない段階で有効化するとエラーになるため、
最初はHTTP専用の設定を作成します。

```bash
mkdir -p /var/www/letsencrypt/.well-known/acme-challenge
nano /etc/nginx/sites-available/space-now
```

以下を保存します。

```nginx
server {
    listen 80;
    server_name 143.198.80.71;

    location ^~ /.well-known/acme-challenge/ {
        root /var/www/letsencrypt;
        default_type text/plain;
    }

    location / {
        return 404;
    }
}
```

設定を有効にし、検証します。

```bash
ln -s /etc/nginx/sites-available/space-now /etc/nginx/sites-enabled/space-now
nginx -t
```

検証が成功した場合だけ、起動・反映します。

```bash
systemctl enable --now nginx
systemctl reload nginx
```

ファイアウォールを使用する場合は、SSH接続を維持したうえで、
外部からTCPポート80と443へアクセスできるよう設定します。
ポート80は証明書更新にも必要です。

確認用ファイルを作ります。

```bash
echo "https-check-ok" > /var/www/letsencrypt/.well-known/acme-challenge/check
```

手元のブラウザで以下を開き、
`https-check-ok` と表示されることを確認します。

http://143.198.80.71/.well-known/acme-challenge/check

### 5. Certbotをインストール

```bash
apt install -y snapd
snap install --classic certbot
/snap/bin/certbot --version
```

IPアドレスのwebroot認証に対応するCertbot 5.4以上を使用します。

### 6. 証明書を取得

```bash
/snap/bin/certbot certonly \
  --webroot \
  --webroot-path /var/www/letsencrypt \
  --preferred-profile shortlived \
  --ip-address 143.198.80.71
```

対話画面でメールアドレスを入力し、利用規約を確認します。

成功すると、以下にファイルが保存されます。

- 証明書：`/etc/letsencrypt/live/143.198.80.71/fullchain.pem`
- 秘密鍵：`/etc/letsencrypt/live/143.198.80.71/privkey.pem`

### 7. HTTPS設定を反映

証明書の取得成功後に実行します。

```bash
cd /opt/space-now
cp deploy/nginx/space-now.conf /etc/nginx/sites-available/space-now
nginx -t
```

検証が成功した場合だけ反映します。

```bash
systemctl reload nginx
```

### 8. 自動更新後の反映を設定

```bash
mkdir -p /etc/letsencrypt/renewal-hooks/deploy

install -m 755 deploy/letsencrypt/reload-nginx \
  /etc/letsencrypt/renewal-hooks/deploy/reload-nginx

/snap/bin/certbot renew --dry-run --run-deploy-hooks
systemctl list-timers --all | grep certbot
```

更新テストが成功し、自動更新用タイマーが表示されることを確認します。

IPアドレス用証明書の有効期間は約6日間なので、
自動更新と更新後のNginx再読み込みを維持してください。

## HTML・CSS・JS・ドキュメントの更新

```bash
cd /opt/space-now
git pull --ff-only origin chore/deployment
```

これらのファイルだけの変更なら、
Goの再ビルドやサービスの再起動は不要です。
ブラウザで再読み込みして確認してください。

## Goコードの更新

まず新しい実行ファイルを作ります。

```bash
cd /opt/space-now
git pull --ff-only origin chore/deployment
/usr/local/go/bin/go build -o /opt/space-now/space-now.new .
```

ビルドが成功した場合だけ、入れ替えて再起動します。

```bash
mv /opt/space-now/space-now.new /opt/space-now/space-now
systemctl restart space-now
systemctl status space-now --no-pager
```

## Nginx設定の更新

```bash
cd /opt/space-now
git pull --ff-only origin chore/deployment

cp /etc/nginx/sites-available/space-now /etc/nginx/space-now.previous.conf
cp deploy/nginx/space-now.conf /etc/nginx/sites-available/space-now
nginx -t
```

検証が成功した場合だけ反映します。

```bash
systemctl reload nginx
```

検証に失敗した場合は再読み込みせず、設定を修正するか、
バックアップを戻してください。

```bash
cp /etc/nginx/space-now.previous.conf /etc/nginx/sites-available/space-now
nginx -t
```

## 動作確認

```bash
systemctl is-active nginx
systemctl is-active space-now
curl -i http://127.0.0.1:8000/api/health
curl -i https://143.198.80.71/api/health
```

両サービスが `active`、healthが200と
`{"status":"OK"}` を返すことを確認します。

HTTPからHTTPSへの転送も確認します。

```bash
curl -I http://143.198.80.71/web/
```

301と、HTTPSへの `Location` ヘッダーが返ることを確認します。

ブラウザでも以下を確認します。

- 証明書の警告なくHTTPSで開ける
- 地図・都市選択・ISS現在位置が表示される
- デモとAPIドキュメントを往復できる

## エラーの確認

```bash
journalctl -u space-now -n 50 --no-pager
tail -n 50 /var/log/nginx/error.log
tail -n 50 /var/log/letsencrypt/letsencrypt.log
```

## 今後の対応

- 詳細API完成後の通過一覧・詳細の実API接続
- mainへのマージ後、配置ブランチと本手順の更新