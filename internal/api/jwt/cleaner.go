package jwt

import (
	"context"
	"cook_book/internal/services"
	"cook_book/pkg/logger"
	"sync"
	"time"
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

// горутина с очистителем
func (s *CleanupService) StartCleanup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancelFunc = cancel

	// задаем интервал - каждые сутки
	ticker := time.NewTicker(24 * time.Hour)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer ticker.Stop()
		for {
			select {
			// когда срабатывает тикер - удаляем
			case <-ticker.C:
				s.cleanExpiredTokens(ctx)
			// когда сворачивается контекст - останавливаем таймер и завершаем клинер
			case <-ctx.Done():
				s.logger.Info("Cleanup scheduler stopped")
				return
			}
		}
	}()
}

// сам очистиетель
func (s *CleanupService) cleanExpiredTokens(ctx context.Context) {
	// создаем контекст с таймаутом
	cleanupCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := s.userUseCase.CleanExpiredTokens(cleanupCtx)
	if err != nil {
		s.logger.Error("Failed to clean expired tokens", "error", err)
		return
	}

	s.logger.Info("Expired tokens cleanup completed successfully")
}

// грейсфул остановка очистителя
func (s *CleanupService) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	s.wg.Wait()
	s.logger.Info("Cleanup service stopped gracefully")
}
