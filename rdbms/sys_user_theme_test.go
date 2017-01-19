package rdbms

import (
	"fmt"
	"testing"
)

func TestSysUserTheme(t *testing.T) {
	//fmt.Println(AddSysUserTheme("zhan", "1001"))
	fmt.Println(GetSysUserTheme("zhan"))
	fmt.Println(UpdateSysUserTheme("zhan", "1001"))
	fmt.Println(GetSysUserTheme("zhan"))
	//fmt.Println(DeleteSysUserTheme("zhan"))
}
