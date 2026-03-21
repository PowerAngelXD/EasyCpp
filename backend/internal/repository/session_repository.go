package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Session 表示保存在 Redis 中的登录会话信息。
type Session struct {
	SessionID  string `json:"sessionId"`            // 会话唯一标识。
	UserID     uint64 `json:"userId"`               // 会话归属用户 ID。
	UserAgent  string `json:"userAgent,omitempty"`  // 客户端标识信息。
	RemoteAddr string `json:"remoteAddr,omitempty"` // 客户端来源地址。
	IssuedAt   int64  `json:"issuedAt"`             // 会话签发时间戳（Unix 秒）。
	ExpiresAt  int64  `json:"expiresAt"`            // 会话过期时间戳（Unix 秒）。
}

type SessionRepository interface {
	Save(ctx context.Context, session Session, ttl time.Duration) error
	Get(ctx context.Context, sessionID string) (*Session, error)
	Delete(ctx context.Context, sessionID string) error
	ListByUserID(ctx context.Context, userID uint64) ([]Session, error)
}

// RedisSessionRepository 提供基于 Redis 的会话存储实现。
type RedisSessionRepository struct {
	redis *redis.Client // Redis 客户端实例。
}

func NewRedisSessionRepository(redisClient *redis.Client) *RedisSessionRepository {
	return &RedisSessionRepository{redis: redisClient}
}

func (r *RedisSessionRepository) Save(ctx context.Context, session Session, ttl time.Duration) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()
	pipe.Set(ctx, sessionKey(session.SessionID), payload, ttl)
	pipe.SAdd(ctx, userSessionSetKey(session.UserID), session.SessionID)
	pipe.Expire(ctx, userSessionSetKey(session.UserID), ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisSessionRepository) Get(ctx context.Context, sessionID string) (*Session, error) {
	payload, err := r.redis.Get(ctx, sessionKey(sessionID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	var session Session
	if err := json.Unmarshal(payload, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *RedisSessionRepository) Delete(ctx context.Context, sessionID string) error {
	session, err := r.Get(ctx, sessionID)
	if err != nil {
		if err == ErrUnauthorized {
			return nil
		}
		return err
	}
	pipe := r.redis.TxPipeline()
	pipe.Del(ctx, sessionKey(sessionID))
	pipe.SRem(ctx, userSessionSetKey(session.UserID), sessionID)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisSessionRepository) ListByUserID(ctx context.Context, userID uint64) ([]Session, error) {
	sessionIDs, err := r.redis.SMembers(ctx, userSessionSetKey(userID)).Result()
	if err != nil {
		return nil, err
	}
	if len(sessionIDs) == 0 {
		return []Session{}, nil
	}

	sessions := make([]Session, 0, len(sessionIDs))
	for _, id := range sessionIDs {
		session, err := r.Get(ctx, id)
		if err != nil {
			continue
		}
		sessions = append(sessions, *session)
	}
	return sessions, nil
}

func sessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

func userSessionSetKey(userID uint64) string {
	return fmt.Sprintf("user_sessions:%d", userID)
}
