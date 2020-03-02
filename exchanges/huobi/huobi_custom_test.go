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

func TestGetContractAccountPositionInfo(t *testing.T) {
	t.Parallel()

	accountPositionInfoData, err := h.GetContractAccountPositionInfo(SymbolBaseType{
		Symbol: "BTC",
	})
	if err != nil {
		t.Errorf("Huobi TestGetContractAccountPositionInfo: %s", err)
	}
	for _, d := range accountPositionInfoData {
		fmt.Printf(" - 帐户信息\n")
		fmt.Printf("\t 合约品种 :%s\n", d.Symbol)
		fmt.Printf("\t 账户权益 :%.8f\n", d.MarginBalance)
		fmt.Printf("\t 持仓保证金 :%.8f\n", d.MarginPosition)
		fmt.Printf("\t 冻结保证金 :%.8f\n", d.MarginFrozen)
		fmt.Printf("\t 可用保证金 :%.8f\n", d.MarginAvailable)
		fmt.Printf("\t 已实现盈亏 :%.8f\n", d.ProfitReal)
		fmt.Printf("\t 未实现盈亏 :%.8f\n", d.ProfitUnReal)
		fmt.Printf("\t 保证金率 :%.8f\n", d.RiskRate)
		fmt.Printf("\t 可划转数量 :%.8f\n", d.WithdrawAvailable)
		fmt.Printf("\t 预估爆仓价 :%.8f\n", d.LiquidationPrice)
		fmt.Printf("\t 杠杆倍数 :%.8f\n", d.LeverRate)
		fmt.Printf("\t 调整系数 :%.8f\n", d.AdjustFactor)
		fmt.Printf("\t 静态权益 :%.8f\n", d.MarginStatic)
		for _, p := range d.Positions {
			fmt.Printf("\t - 持仓信息\n")
			fmt.Printf("\t\t 合约类型 :%s\n", p.ContractType)
			fmt.Printf("\t\t 合约代码 :%s\n", p.ContractCode)
			fmt.Printf("\t\t 持仓量 :%.8f\n", p.Volume)
			fmt.Printf("\t\t 可平仓数量 :%.8f\n", p.Available)
			fmt.Printf("\t\t 冻结数量 :%.8f\n", p.Frozen)
			fmt.Printf("\t\t 开仓均价 :%.8f\n", p.CostOpen)
			fmt.Printf("\t\t 持仓均价 :%.8f\n", p.CostHold)
			fmt.Printf("\t\t 未实现盈亏 :%.8f\n", p.ProfitUnreal)
			fmt.Printf("\t\t 收益率 :%.8f\n", p.PofitRate)
			fmt.Printf("\t\t 收益 :%.8f\n", p.Pofit)
			fmt.Printf("\t\t 持仓保证金 :%.8f\n", p.PositionMargin)
			fmt.Printf("\t\t 杠杆倍数 :%d\n", p.LeverRate)
			fmt.Printf("\t\t 交易方向 :%s\n", p.Direction)
			fmt.Printf("\t\t 最新价 :%.8f\n", p.LastPrice)
		}
	}
}
