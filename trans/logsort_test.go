package trans

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestCsvWrite(t *testing.T) {
	t1 := Transaction{1, 1, 1, 1, 1}
	t2 := Transaction{2, 4, 1, 1, 1}
	t3 := Transaction{3, 3, 1, 1, 1}

	trs := Transactions{t1, t2, t3}

	CsvWrite(trs, os.Stdout)

	trs.time_sort()
	CsvWrite(trs, os.Stdout)
}

func TestLogLoad(t *testing.T) {
	r := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	r.time_sort()
}

func TestLogLoad2(t *testing.T) {
	tran := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	tran.time_sort()

	var max_rec int

	for i := range tran {
		action := tran[i].Action
		if action == PARTIAL || action == UPDATE_BUY || action == UPDATE_SELL {
			continue
		}
		t := tran[i]
		fmt.Printf("%2d %20d %10d %20d %16d\n", t.Action, t.Time_stamp, t.Price, t.Volume, t.NextTime)

		if 200 < max_rec {
			break
		}
		max_rec += 1
	}
}

func TestLogLoad3(t *testing.T) {
	tran := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	tran.time_sort()

	var max_rec int

	var trans Transactions
	var last_oi int64
	var diff_oi int64

	for i := range tran {
		action := tran[i].Action
		if action == PARTIAL || action == UPDATE_BUY || action == UPDATE_SELL {
			continue
		}
		// t := tran[i]

		if action == TRADE_SELL || action == TRADE_BUY {
			trans = append(trans, tran[i])
		}

		if action == OPEN_INTEREST {
			diff_oi = tran[i].Volume - last_oi
			fmt.Println(trans)
			fmt.Println(diff_oi)

			if len(trans) != 0 {
				// Dp2(trans, 0, diff_oi)
				//mask := FindBestMatchMask(trans, int(diff_oi))
				//fmt.Println(mask)
			}

			trans = make(Transactions, 0)

			last_oi = tran[i].Volume
		}

		// fmt.Printf("%2d %20d %10d %20d %16d\n", t.Action, t.Time_stamp, t.Price, t.Volume, t.NextTime)

		/*
			if 200 < max_rec {
				break
			}
		*/
		max_rec += 1
	}
}

// {x1, x2, x3, x4} -> {x1, x1, x2, x2, x3, x3, x4, x4}
//
//
//

var result_item []int

/*

func Dp(item []int, max_value int) {
	n := len(item)

	dp := make([][]int64, n+1)

	for i := 0; i < n+1; i++ {
	dp[i] = make([]int64, M+1)

}
// 初期化
for i := 0; i <= N; i++ {
	for j := 1; j <= M; j++ {
		dp[i][j] = -1
	}
}
dp[0][0] = 0
// traceにはどこの(i, j)から来たかの情報を保持させる
for i := 1; i <= N; i++ {
	for j := 0; j <= M; j++ {
		if dp[i-1][j] >= 0 {
			dp[i][j] = dp[i-1][j]
			trace[i][j] = Route{i - 1, j}
		}
		if j >= weights[i-1] && dp[i-1][j-weights[i-1]] >= 0 && dp[i][j] < dp[i-1][j-weights[i-1]]+values[i-1] {
			dp[i][j] = dp[i-1][j-weights[i-1]] + values[i-1]
			trace[i][j] = Route{i - 1, j - weights[i-1]}
		}
	}
}
// 出力
fmt.Println(dp[N][M])
// DPテーブルの確認
fmt.Println(dp)
pnt := Route{N, M}
routes := []Route{}
res := []int{}
// トレースバック
for pnt.x != 0 {
	routes = append(routes, pnt)
	pre := pnt
	pnt = trace[pnt.x][pnt.y]
	// 異なる重さのところから値が来ている場合は荷物iを追加したとき
	if pre.y != pnt.y {
		res = append(res, pre.x)
	}
}
fmt.Println(routes)
// 重さがWになるときのベストな荷物の組み合わせの出力
fmt.Println(res)
}

}




*/

func Dp(item []int, i int, max_vol int) (result int, weight int) {
	fmt.Println("index=", i, "maxvol=", max_vol)

	if i == -1 {
		fmt.Println("ERR")
		return -1, 0
	} else {
		w := item[i]

		if max_vol < w {
			fmt.Print("-", w, "-")
			return Dp(item, i-1, max_vol)
		} else if w < max_vol {
			a, wa := Dp(item, i-1, max_vol-w)
			b, wb := Dp(item, i-1, max_vol)
			if a < b {
				fmt.Print("-", w, "-")
				return b, wb
			} else {
				fmt.Print("[", w, "] ")
				return a, wa
			}
		} else if w == max_vol {
			fmt.Println("----HIT-----", w)
			return 0, 0
		}
	}

	log.Println("unknown status")
	return 0, 0 // err
}
