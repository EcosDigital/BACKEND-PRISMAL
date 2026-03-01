package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ecosistema/core/src/core"
	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	EmpresaID *int      `json:"empresa_id,omitempty"`
	SedeID    *int      `json:"sede_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// guardar session en redis
func SaveSession(jti string, session SessionData, ttl time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", jti)

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return core.RedisClient.Set(ctx, key, data, ttl).Err()
}

// validar si la session existe en redis
func ValidateSession(jti string) (*SessionData, error) {
	if jti == "" {
		return nil, fmt.Errorf("session invalid")
	}

	ctx := context.Background()
	key := fmt.Sprintf("session:%s", jti)

	data, err := core.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session invalid")
	}

	if err != nil {
		return nil, err
	}

	var session SessionData
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil

}

// revocar session individual
func RevokeSession(jti string) error {
	if jti == "" {
		return fmt.Errorf("jti vacío")
	}

	ctx := context.Background()
	key := fmt.Sprintf("session:%s", jti)

	result, err := core.RedisClient.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	// Verificar que realmente se eliminó
	if result == 0 {
		return fmt.Errorf("sesión no encontrada en Redis")
	}

	return nil
}

// revocar todas las sessiones de un usuario
func RevokeAllUserSessions(userID int) error {
	ctx := context.Background()
	pattern := "session:*"

	var cursor uint64
	var deletedCount int

	for {
		keys, nextCursor, err := core.RedisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			data, err := core.RedisClient.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var session SessionData
			if err := json.Unmarshal([]byte(data), &session); err != nil {
				continue
			}

			if session.UserID == userID {
				core.RedisClient.Del(ctx, key)
				deletedCount++
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
