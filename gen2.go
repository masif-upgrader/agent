// +build ignore

package main

import (
	"fmt"
	"github.com/masif-upgrader/common"
	"os"
)

func main() {
	GithubcomAl2klimovGogeneratedeps["github.com/masif-upgrader/agent"] = struct{}{}
	GithubcomAl2klimovGogeneratedeps["github.com/golang/go"] = struct{}{}

	if errGC := common.GenCredits("main", GithubcomAl2klimovGogeneratedeps); errGC != nil {
		fmt.Println(errGC.Error())
		os.Exit(1)
	}
}
