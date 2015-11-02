package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"time"

	"github.com/tpjg/goriakpbc"
)

// Logger パッケージ内で利用出来るよう定義
var Logger *log.Logger

func main() {
	// error 変数定義
	var err error
	// ip address
	var flagAddr = flag.String("i", "127.0.0.1:8087", "Riak Server IP Address and Port")
	var flagProtocol = flag.String("p", "TCP", "Protocol (TCP/PB)")
	var flagVersion = flag.Bool("v", false, "Print Version Info")
	// 引数のparse
	flag.Parse()

	if *flagVersion {
		fmt.Printf("riak-ping version: %s\n", Version)
		os.Exit(0)
	}

	// log書き込み定義
	Logger, err = SetLogger()
	if err != nil {
		panic(err)
	}

	// riak clientの定義
	myRiak := SetRiakClient(*flagAddr)

	// riak:疎通確認
	err = CheckConnect(myRiak, *flagProtocol)
	if err != nil {
		os.Exit(1)
	}

	return
}

// SetLogger log書き込み設定の定義
func SetLogger() (*log.Logger, error) {
	// Logger を crit / syslog に設定
	l, err := syslog.NewLogger(syslog.LOG_CRIT|syslog.LOG_SYSLOG, 0)
	return l, err
}

// SetRiakClient riak clientの定義および初期値設定
func SetRiakClient(a string) *riak.Client {
	// Client の定義
	c := riak.NewClient(a)
	return c
}

// CheckConnect 疎通確認
func CheckConnect(c *riak.Client, t string) error {
	// ステータスchannel
	resultCh := make(chan error, 1)
	// 疎通確認は10secで切断したいので非同期で実施
	go ConnectRiak(c, t, resultCh)

	// 疎通確認完了 or 10sec経てば CheckConnect を終了
	select {
	case err := <-resultCh: // 疎通確認が完了(終了フラグがtrue)
		if err != nil {
			// 疎通NG: 画面とsyslogに記載
			fmt.Printf("%s Connect Failure.\n", t)
			if Logger != nil {
				Logger.Printf("%s Connect Failure.\n", t)
			}
		} else {
			// 疎通OK: 画面のみに記載
			fmt.Printf("%s Connect Success.\n", t)
		}
		return err
	case <-time.After(time.Second * 10): // 疎通確認処理が終わらず10秒経過した場合
		// 疎通確認自体がNG: 画面とsyslogに記載
		fmt.Printf("%s Connect Timed Out.\n", t)
		if Logger != nil {
			Logger.Printf("%s Connect Timed Out.\n", t)
		}
		// 戻り値のerrorを作成
		err := errors.New("timed out")
		return err
	}
}

// ConnectRiak riakへの接続選択(TCP/PB)および疎通結果の取得
func ConnectRiak(c *riak.Client, t string, rch chan error) {
	var err error
	defer c.Close()
	// 接続方式の選択
	switch {
	case t == "TCP":
		// L4レベルでの疎通確認
		err = c.Connect()
	case t == "PB":
		// L7レベルでの疎通確認
		err = c.Ping()
	default:
		err = errors.New("unknown protocol")
	}
	// 接続結果を error Channel に返す (error型がnilの場合rchに直接渡せないのでifで分岐している)
	if err != nil {
		rch <- err
	} else {
		// rch <- err とした場合、err=nilだとロックされてしまうのでnilを直接渡してる。
		rch <- nil
	}
}
