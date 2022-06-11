// 並行処理で url,HTTP_Status の状態で出力
// STDIN ではURLの読み込み
// 第一引数を同時並行リクエスト数の上限と認識
//	(キャップを掛けないとソケットの 数が不足してエラーになる)
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const CONCURR_LIMIT_DEFAULT = 5

func main() {
	var concurr_limit int
	if len(os.Args) > 1 {
		var err error
		concurr_limit, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		concurr_limit = CONCURR_LIMIT_DEFAULT
	}
	var wg sync.WaitGroup
	ch := make(chan string, concurr_limit) // limited buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			defer func() { fmt.Println(<-ch) }() // 処理後チャンネルを解放
			resp, err := http.Get(url)
			if err != nil {
				ch <- fmt.Sprintf("%s,%v", url, err)
				return
			}
			ch <- fmt.Sprintf("%s,%v", url, resp.Status)
			return
		}(url)
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
