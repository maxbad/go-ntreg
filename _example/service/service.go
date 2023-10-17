package main

import (
	"fmt"
	"github.com/maxbad/go-ntreg"
	"os"
)

func main() {
	fileName, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取当前目录：", err)
		return
	}
	fileName += "/SYSTEM"
	serviceName := "xadTest001"
	fmt.Printf("serviceName:%s fileName:%s \n", serviceName, fileName)

	if true {
		if e := ntreg.CreateService(serviceName, fileName); e != nil {
			fmt.Printf("[%s]创建失败! %s \n", serviceName, e.Error())
			return
		}
		fmt.Printf("[%s]创建成功! \n", serviceName)
	}

	if true {
		if e := ntreg.DeleteService(serviceName, fileName); e != nil {
			fmt.Printf("[%s]删除失败! %s \n", serviceName, e.Error())
			return
		}
		fmt.Printf("[%s]删除成功! \n", serviceName)
	}

}
