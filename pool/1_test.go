package pool

import (
	"fmt"
	"testing"
)

func Test_initPool(t *testing.T) {
	initPool()

	p := pool.Get().(*Person)
	fmt.Println("首次从 pool 里获取：", p)

	p.Name = "first"
	fmt.Printf("设置 p.Name = %s\n", p.Name)

	pool.Put(p)

	fmt.Println("Pool 里已有一个对象：&{first}，调用 Get: ", pool.Get().(*Person))
	temp := pool.Get().(*Person)
	fmt.Println("Pool 没有对象了，调用 Get,返回新的对象", &temp)

	p.Name = "Second"
	pool.Put(p)
	fmt.Println("Pool 又有一个对象，调用 Get: ", pool.Get().(*Person))
}

func BenchmarkTest1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		temp := &Person{
			Name: "aaa",
		}
		_ = temp.Name
	}
}

func BenchmarkTest2(b *testing.B) {
	initPool()
	for i := 0; i < b.N; i++ {
		temp := pool.Get().(*Person)
		_ = temp.Name
		pool.Put(temp)
	}
}
