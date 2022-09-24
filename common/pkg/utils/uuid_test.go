package utils_test

import (
	"github.com/tmnhs/crony/common/pkg/utils"
	"testing"
)

func TestUUID(t *testing.T) {
	uuid, err := utils.UUID()
	if err != nil {
		t.Error(err)
	}
	t.Log(uuid)
}
