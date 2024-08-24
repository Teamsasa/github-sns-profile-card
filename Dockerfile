# Goの公式イメージをベースにする
FROM golang:latest

# 作業ディレクトリを設定
WORKDIR /app

# Goのモジュールを有効にする
ENV GO111MODULE=on

# 必要なGoのパッケージをダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY ./ ./

# APIサーバーをビルド
RUN make build

# 環境変数を設定
ENV PORT=8080

# ポート8080を公開
EXPOSE 8080

# スクリプトを実行
CMD ["/app/main"]