package huobi

import (
	"fmt"
	"testing"
)

func TestGetContractAccountInfo(t *testing.T) {
	t.Parallel()

	accountInfo, err := h.GetContractAccountInfo(ContractAccountInfoRequest{})
	if err != nil {
		t.Errorf("Huobi TestGetContractAccountInfo: %s", err)
	} else {
		for _, v := range accountInfo {
			fmt.Println(v)
		}
	}
}
