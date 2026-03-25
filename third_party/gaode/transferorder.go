package gaode

type TransferOrder struct {
	OrderID string   // 全局唯一业务单号
	BatchID string   // 批次号
	UserID  string   
	AccountType int   // 支付宝/微信/银行卡
	AccountNo  string // 账号
	Amount     int64  // 单位: 分
	Status     int    // 待审核/待下发/处理中/成功/失败/退票
	FreezeID   string // 冻结记录
	RiskResult int    // 风控结果
	AuditStatus int   // 待审核/通过/拒绝
	CreateAt    time.Time // 创建时间
}

// 可用余额 冻结余额 在途余额 已结算金额
// 检查可用余额 >= 总下发金额
// 生成冻结账单 ， 冻结对应金额
// 冻结成功 -> 生成代发单
// 下发成功 -> 扣减成功 -> 正式记账
// 下发失败 / 退票 -> 解冻退回可用余额
