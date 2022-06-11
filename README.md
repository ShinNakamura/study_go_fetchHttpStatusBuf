# Golang: Fetch HTTP Status with buffer for cap

複数のURLに対して並行でHTTP Statusを取得する。ただし、並行処理が同時に走る数をバッファで制御する。: Golang 自習課題

## コマンドライン引数の数で同時処理数を制限するバージョン

過去に作った元になるバージョンはこれ。

Goのコード内ではなく、Goプログラムを起動するコマンドライン上でxargsを使って引数をコントロールし
並行処理の対象になる全体数をそもそも制限している。

```go
// 並行処理で url,HTTP_Status の状態で出力
// STDIN ではURLの読み込み数＝リクエスト数を制限するのが逆に煩雑になる(キャップを掛けないとソケットの 数が不足してエラーになる)
//  なので os.Args でURLを受け取りるようにして、URL数は xargs -L n で制御する
package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    ch := make(chan string)
    for _, url := range os.Args[1:] {
        go fetch(url, ch) // ゴルーチン開始
    }
    for range os.Args[1:] {
        fmt.Println(<-ch)
    }
}

func fetch(url string, ch chan<- string) {
    resp, err := http.Get(url)
    if err != nil {
        ch <- fmt.Sprintf("%s,%v", url, err)
        return
    }
    ch <- fmt.Sprintf("%s,%v", url, resp.Status)
    return
}
```

## 今回作る改良版

- コマンドラインでの使用では、単純にURLのリストを標準入力から与えるようにしたい(＝xargsを使わない。)
- 並行処理の同時実行数上限を第一引数で与えることができるようにしたい。(指定がなければデフォルト数を適用)


