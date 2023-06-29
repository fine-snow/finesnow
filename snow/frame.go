// Framework information

package snow

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	ver    = "v0.0.1 beta"
	oldStr = "snow/frame.go"
	newStr = "banner.txt"
)

// outputFrameworkInfo Output information such as framework logo and version
func outputFrameworkInfo() {
	// Read the banner.text
	_, f, _, _ := runtime.Caller(0)
	url := strings.Replace(f, oldStr, newStr, -1)
	file, err := os.Open(url)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	if err = scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println(ver)
	_ = file.Close()
}
