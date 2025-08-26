package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"worker-manager/internal/domain"
	"worker-manager/util/logger"

	"github.com/valkey-io/valkey-go"
)

type ValkeyRepository struct {
	client valkey.Client
}

var _ domain.DataRepository = (*ValkeyRepository)(nil)

func NewValkeyRepository(valkeyAddress string) *ValkeyRepository {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkeyAddress},
	})
	if err != nil {
		logger.Errorln(context.Background(), "Failed to connect to Valkey at", valkeyAddress, ":", err)
		panic(err)
	}
	logger.Infoln(context.Background(), "Connected to Valkey at", valkeyAddress)
	return &ValkeyRepository{
		client: client,
	}
}

func (r *ValkeyRepository) Put(ctx context.Context, key domain.WorkerId, value *domain.Data) error {
	strData := toJson(value)
	cmd := r.client.B().Set().Key(strconv.Itoa(int(key))).Value(strData).Build()
	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		logger.Errorf(ctx, "Failed to put value in valkey: %v, %v, %v", key, strData, err)
		return fmt.Errorf("failed to put value in valkey: %w", err)
	}
	logger.Debugf(ctx, "put value in valkey: %v, %v", key, strData)
	return nil
}

func (r *ValkeyRepository) Get(ctx context.Context, key domain.WorkerId) (*domain.Data, error) {
	cmd := r.client.B().Get().Key(strconv.Itoa(int(key))).Build()
	strData, err := r.client.Do(ctx, cmd).ToString()
	if err != nil {
		logger.Errorf(ctx, "Failed to get value in valkey: %v, %v", key, err)
		return nil, fmt.Errorf("failed to get value from valkey: %w", err)
	}
	data, err := fromJson(strData)
	if err != nil {
		logger.Errorf(ctx, "Failed to convert json value from valkey: %v, %v", key, err)
		return nil, fmt.Errorf("failed to parse data from valkey: %w", err)
	}
	logger.Debugf(ctx, "get value from valkey: %v, %v", key, strData)
	return data, nil
}

func (r *ValkeyRepository) Delete(ctx context.Context, key domain.WorkerId) error {
	cmd := r.client.B().Del().Key(strconv.Itoa(int(key))).Build()
	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		logger.Errorf(ctx, "Failed to delete value from valkey: %v, %v", key, err)
		return fmt.Errorf("failed to delete value from valkey: %w", err)
	}
	logger.Debugf(ctx, "deleted value from valkey: %v", key)
	return nil
}

func toJson(d *domain.Data) string {
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err).Error()
	}
	return string(b)
}

func fromJson(data string) (*domain.Data, error) {
	var d domain.Data
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return &d, nil
}
