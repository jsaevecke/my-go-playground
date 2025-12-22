package migration

import "context"

type Service struct {
	migrator Migrator
}

func New(migrator Migrator) *Service {
	return &Service{
		migrator: migrator,
	}
}

func (svc *Service) Up(ctx context.Context) error {
	return svc.migrator.Up(ctx)
}

func (svc *Service) Status(ctx context.Context) error {
	return svc.migrator.Status(ctx)
}
