package datastruct

// SyncKeyItem 同步Key的元素
type SyncKeyItem struct {
	Key int64 `json:"Key"`
	Val int64 `json:"Val"`
}

// SyncKey 同步Key
type SyncKey struct {
	Count int64          `json:"Count"`
	List  []*SyncKeyItem `json:"List"`
}
