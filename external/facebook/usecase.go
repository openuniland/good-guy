package facebook

import "context"

type UseCase interface {
	SendMessage(ctx context.Context, id string, message interface{}) error
}
