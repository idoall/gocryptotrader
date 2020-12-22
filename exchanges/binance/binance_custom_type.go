package binance

import (
	"time"

	"github.com/idoall/gocryptotrader/exchanges/asset"
)

// AccountSnapshotRequest 查询每日资产快照 (USER_DATA)
type AccountSnapshotRequest struct {
	Type      asset.Item `json:"type"`
	Price     float64    `json:"price"`
	Limit     int64      `json:"limit"`
	StartTime int64      `json:"startTime"`
	EndTime   int64      `json:"endTime"`
}

// AccountSnapshotResponse 查询每日资产快照 (USER_DATA) - 返回信息
type AccountSnapshotResponse struct {
	TotalAssetOfBtc float64    `json:"totalAssetOfBtc"`
	Asset           asset.Item `json:"asset"`
	Symbol          string     `json:"symbol"`
	Free            float64    `json:"free"`
	Locked          float64    `json:"locked"`
	UpdateTime      time.Time  `json:"updateTime"`
}
