package store

import (
	"context"
	"fmt"
	"time"

	"github.com/lesi97/go-av-scanner/internal/utils"
)

func (s *DbApiStore) Health(ctx context.Context) (*string, error) {
	context, ok := ctx.Value(ContextKey).(Context)
	if !ok {
		return nil, fmt.Errorf("failed to get context value")
	}

	utils.PrintPrettyJSON(context)
	msg := time.Now().UTC().Format(time.RFC3339)
	return &msg, nil
}