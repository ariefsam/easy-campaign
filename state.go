package campaign

import (
	"context"
)

type Step func(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error)

type InternalState struct {
	Session struct {
		UserID string
	}
}
