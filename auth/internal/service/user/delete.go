package user

import (
	"auth/internal/logger"
	"auth/internal/utils/kafka_common"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	logger.Info("Deleting user...", zap.Int64("id", id))
	_, user, _, _, err := s.userRepository.Get(ctx, id)
	if err != nil {
		logger.Error("Failed to get user for deleting", zap.Int64("id", id), zap.Error(err))
		return fmt.Errorf("failed to get user for deleting: %w", err)
	}

	// TODO: транзакция для удаления и отправки сообщения в кафку
	err = s.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	logger.Info("User deleted",
		zap.Int64("id", id),
		zap.String("name", user.Name),
		zap.Int32("role", user.Role))

	deletionTime := time.Now()
	event := kafka_common.UserDeletedPayload{UserID: id}
	dataEvent, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal event",
			zap.Error(err),
			zap.ByteString("event", dataEvent),
			zap.Time("event_time", deletionTime),
			zap.String("topic", kafka_common.UserDeletedTopic),
		)
		return fmt.Errorf("Failed to marshal event: %w", err)
	}

	err = s.producer.Send(ctx, nil, dataEvent, kafka_common.UserDeletedTopic, deletionTime)
	if err != nil {
		logger.Error("Failed to send event",
			zap.Error(err),
			zap.ByteString("event", dataEvent),
			zap.Time("event_time", deletionTime),
			zap.String("topic", kafka_common.UserDeletedTopic),
		)
		return fmt.Errorf("Failed to send event: %w", err)
	}

	return nil
}
