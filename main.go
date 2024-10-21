package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
)

func main() {
	generateCertAndKey()

	s := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &customHandler{},
	}

	log.Println(http2.ConfigureServer(&s, &http2.Server{}))
	log.Println(s.ListenAndServeTLS("cert.pem", "key.pem"))
	// cert.pemはCAが作成した証明書ではないので、普通にアクセスするとエラーになる
}

// 自己署名のSSL/TLS証明書とRSA秘密鍵を生成
// 本番では認証局から取得する
func generateCertAndKey() {

	// 大きな乱数を生成して、証明書のシリアル番号として使用
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)

	// 証明書の発行者および発行された者の情報を設定
	subject := pkix.Name{
		Organization:       []string{"Hello World Co."}, // 組織名
		OrganizationalUnit: []string{"Customer"},        // 部署名
		CommonName:         "Go Web Programming",        // 一般名(ドメイン名などが通常入る)
	}

	// 証明書テンプレートの作成
	template := x509.Certificate{
		SerialNumber: serialNumber,                         // シリアル番号
		Subject:      subject,                              // 発行者情報
		NotBefore:    time.Now(),                           // 証明書の有効開始日時
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 有効期間は1年間
		KeyUsage: x509.KeyUsageDataEncipherment | // 証明書の使用目的：データ暗号化とデジタル署名
			x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, // クライアント認証用の証明書
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},             // 証明書が対応するIPアドレス（ここではlocalhost）
	}

	// 2048ビットのRSA鍵ペアを生成（秘密鍵と公開鍵）
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	// 証明書を生成（自己署名証明書）
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)

	// 証明書を "cert.pem" ファイルに書き出し
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	// 秘密鍵を "key.pem" ファイルに書き出し
	keyOut, _ := os.Create("key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}
