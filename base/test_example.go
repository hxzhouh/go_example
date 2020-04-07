/*
@Time : 2020/4/7 11:46
@Author : ZhouHui2
@File : test_example
@Software: GoLand
*/
package base

var (
	UserMap = make(map[int64]int64)
)

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func Insert(userId, corpId int64) {

}
