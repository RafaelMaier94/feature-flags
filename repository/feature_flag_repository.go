package repository

import (
	"context"

	"github.com/rafaelmaier/featureflags/domain"
)

type FeatureFlagRepository interface{
	Create(ctx context.Context, flag *domain.FeatureFlag) error

	Update(ctx context.Context, flag *domain.FeatureFlag) error

	Delete(ctx context.Context, key string) error

	Get(ctx context.Context, key string) (*domain.FeatureFlag, error)

	List(ctx context.Context) ([]*domain.FeatureFlag, error)
}