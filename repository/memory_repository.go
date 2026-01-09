package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/rafaelmaier/featureflags/domain"
)

type InMemoryRepository struct {
	mu sync.RWMutex
	flags map[string]*domain.FeatureFlag
}

func NewInMemoryRepository() *InMemoryRepository{
	return &InMemoryRepository{
		flags: make(map[string]*domain.FeatureFlag),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, flag *domain.FeatureFlag) error{
	if flag == nil{
		return errors.New("flag cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.flags[flag.Key]; exists{
		return errors.New("feature flag already exists")
	}
	r.flags[flag.Key] = flag
	return nil
}

func (r *InMemoryRepository) Update(ctx context.Context, flag *domain.FeatureFlag) error{
	if flag == nil{
		return errors.New("flag cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.flags[flag.Key]; !exists{
		return errors.New("feature flag not found")
	}

	r.flags[flag.Key] = flag
	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, key string) error{
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.flags[key]; !exists{
		return errors.New("feature flag not found")
	}

	delete(r.flags, key)
	return nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*domain.FeatureFlag, error){
	r.mu.RLock()
	defer r.mu.RUnlock()

	flags := make([]*domain.FeatureFlag, 0, len(r.flags))
	for _, flag := range r.flags{
		flags = append(flags, flag)
	}
	return flags, nil
}

func (r *InMemoryRepository) Get(ctx context.Context, key string) (*domain.FeatureFlag, error){
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    flag, exists := r.flags[key]
    if !exists {
        return nil, errors.New("feature flag not found")
    }
    
    return flag, nil
}