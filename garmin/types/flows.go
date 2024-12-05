package types

type FlowGenericResp[T any] struct {
	Err error
	Ok  *T
}
