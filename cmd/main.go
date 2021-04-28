package main

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/markxp/tail"
)

const mabi = `C:\Nexon\Mabinogi\Tin_log.txt`
const testFile = `readme.txt`

func main() {

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	r, ec := tail.RetryTailf(ctx, testFile, 0)
	go func() {
		for e := range ec {
			fmt.Println(e)
		}
	}()

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}

	fmt.Println("end")

}
