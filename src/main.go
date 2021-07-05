package main

import (
	"fmt"

	workwx "github.com/xen0n/go-workwx" // package workwx
)

// package workwx

func main() {
	// New 构造一个 Workwx 客户端对象，需要提供企业 ID
	wx := workwx.New("ww04d8763e81ec8099")
	wxApp := wx.WithApp("ZEXAjKHsLxkbIoCdnl3aA8buSfzC2AA7i7OTIvphwso", 1000003)
	info, err := wxApp.ListAllDepts()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}
