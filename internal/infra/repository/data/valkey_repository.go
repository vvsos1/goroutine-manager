package data

import (
	"context"
	"fmt"
	"goroutine-manager/internal/domain"
	"strconv"

	"github.com/valkey-io/valkey-go"
)

type ValkeyRepository struct {
	client valkey.Client
}

var _ domain.KeyValueRepository = (*ValkeyRepository)(nil)

func NewValkeyRepository(valkeyAddress string) *ValkeyRepository {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkeyAddress},
	})
	if err != nil {
		// TODO: 에러 처리
		panic(err)
	}
	fmt.Println("Connected to Valkey at", valkeyAddress)
	return &ValkeyRepository{
		client: client,
	}
}

func (r *ValkeyRepository) Put(key domain.GoroutineId, value string) error {
	ctx := context.Background()
	cmd := r.client.B().Set().Key(strconv.Itoa(int(key))).Value(value).Build()
	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to put value in valkey")
	}
	return nil
}

func (r *ValkeyRepository) Get(key domain.GoroutineId) (string, error) {
	ctx := context.Background()
	cmd := r.client.B().Get().Key(strconv.Itoa(int(key))).Build()
	result, err := r.client.Do(ctx, cmd).ToString()
	if err != nil {
		return "", fmt.Errorf("failed to get value from valkey")
	}
	return result, nil
}
