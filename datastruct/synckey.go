package datastruct

type SyncKeyItem struct {
	Key int64 `json:"Key"`
	Val int64 `json:"Val"`
}

type SyncKey struct {
	Count int64          `json:"Count"`
	List  []*SyncKeyItem `json:"List"`
}
