/*
@Time : 2020/4/8 17:06
@Author : ZhouHui2
@File : Slice_test
@Software: GoLand
*/
package base

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	fmt.Printf("%p\n", &a)
	fmt.Printf("%p\n", &b)
	for k, v := range a {
		fmt.Println(&a[k], "  ", v)
	}
	a = append(a, b...)
	fmt.Printf("%p\n", &a)
}
