package escalation

type EscalationPolicy struct {
	Thresholds []int
}

func DefaultPolicy() EscalationPolicy {
	return EscalationPolicy{Thresholds: []int{5, 10, 15}} // seconds for demo
}
