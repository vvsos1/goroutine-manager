package data

import (
	"context"
	"encoding/json"
	"fmt"
	"goroutine-manager/internal/domain"
	"log"
	"strconv"

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
		log.Println("Failed to connect to Valkey at", valkeyAddress, ":", err)
		panic(err)
	}
	log.Println("Connected to Valkey at", valkeyAddress)
	return &ValkeyRepository{
		client: client,
	}
}

func (r *ValkeyRepository) Put(key domain.WorkerId, data *domain.Data) error {
	strData := ToJson(data)
	ctx := context.Background()
	cmd := r.client.B().Set().Key(strconv.Itoa(int(key))).Value(strData).Build()
	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to put value in valkey: %w", err)
	}
	return nil
}

func (r *ValkeyRepository) Get(key domain.WorkerId) (*domain.Data, error) {
	ctx := context.Background()
	cmd := r.client.B().Get().Key(strconv.Itoa(int(key))).Build()
	strData, err := r.client.Do(ctx, cmd).ToString()
	if err != nil {
		return nil, fmt.Errorf("failed to get value from valkey: %w", err)
	}
	data, err := FromJson(strData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data from valkey: %w", err)
	}
	return data, nil
}

func (r *ValkeyRepository) Delete(key domain.WorkerId) error {
	ctx := context.Background()
	cmd := r.client.B().Del().Key(strconv.Itoa(int(key))).Build()
	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to delete value from valkey: %w", err)
	}
	return nil
}

func ToJson(d *domain.Data) string {
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err).Error()
	}
	return string(b)
}

func FromJson(data string) (*domain.Data, error) {
	var d domain.Data
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return &d, nil
}
