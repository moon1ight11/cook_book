package jwt

import (
	"context"
	"sync"
	"time"

	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

type CleanupService struct {
	userUseCase services.UserUseCaseInterface
	logger      logger.Logger
	wg          sync.WaitGroup
	cancelFunc  context.CancelFunc
}

func NewCleanupService(userUseCase services.UserUseCaseInterface, logger logger.Logger) *CleanupService {
	return &CleanupService{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

func (s *CleanupService) StartCleanup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancelFunc = cancel

	ticker := time.NewTicker(24 * time.Hour)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.cleanExpiredTokens(ctx)
			case <-ctx.Done():
				s.logger.Info("Cleanup scheduler stopped")
				return
			}
		}
	}()
}

func (s *CleanupService) cleanExpiredTokens(ctx context.Context) {
	cleanupCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.userUseCase.CleanExpiredTokens(cleanupCtx); err != nil {
		s.logger.Error("Failed to clean expired tokens", "error", err)
		return
	}

	s.logger.Info("Expired tokens cleaned")
}

func (s *CleanupService) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	s.wg.Wait()
	s.logger.Info("Cleanup service stopped")
}
