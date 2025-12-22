package gormdb

import (
	"context"
)

func (g *GormDatabase) Health(ctx context.Context) map[string]string {
	return g.sqlDB.Health(ctx)
}
