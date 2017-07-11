package plugin

import "context"

type InitContext struct {
	Context context.Context
	Config  interface{}
}
