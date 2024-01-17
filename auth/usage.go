package auth

type KeyUsage struct {
	LastIP      string           `json:"last_ip"`
	LastAccess  int64            `json:"last_access"`
	CallCount   int64            `json:"call_count"`
	MethodUsage map[string]int64 `json:"method_usage"`
}
