package types

import (
	"os"

	"github.com/zsmatrix62/garmin-cli/garmin/pkg/helpers"
)

type FlowGenericResp[T any] struct {
	Err error
	Ok  *T
}

func (f FlowGenericResp[T]) ToStdOut() {
	jRes, _ := helpers.JsonString(f)
	_, _ = os.Stdout.Write([]byte(jRes))
}
