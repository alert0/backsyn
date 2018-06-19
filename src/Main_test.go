package main

import (
	"logger"
	"testing"
)

func Test_compress7zip(t *testing.T) {

	err := compress7zip("D:/3/DLLLIB", "d:/aaa.7z")

	if err != nil {

		logger.Info("压缩错误")
	}

	//if !IsPalindrome("detartrated") {
	//	t.Error(`IsPalindrome("detartrated") = false`)
	//}
	//if !IsPalindrome("kayak") {
	//	t.Error(`IsPalindrome("kayak") = false`)
	//}
}
