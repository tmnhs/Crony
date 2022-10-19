package utils_test

import (
	"github.com/tmnhs/crony/common/pkg/utils"
	"testing"
)

func TestMD5(t *testing.T) {
	t.Log(utils.MD5("123456"))
}
