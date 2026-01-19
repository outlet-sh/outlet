package workers

import (
	"github.com/outlet-sh/outlet/internal/services/email"
	"github.com/outlet-sh/outlet/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// EmailWorker wraps the dispatcher for the worker interface
type EmailWorker struct {
	dispatcher *email.Dispatcher
}

// StartEmailWorker starts the high-performance email dispatcher
func StartEmailWorker(svcCtx *svc.ServiceContext) *EmailWorker {
	logx.Info("Initializing email worker...")

	// Create sequence service
	sequenceService := email.NewSequenceServiceWithBaseURL(
		svcCtx.DB,
		svcCtx.EmailService,
		svcCtx.Config.App.BaseURL,
	)

	// Build dispatcher config from service context config
	config := email.DefaultDispatcherConfig()

	// Override with config values if present
	if svcCtx.Config.Email.WorkerCount > 0 {
		config.Workers = svcCtx.Config.Email.WorkerCount
	}
	if svcCtx.Config.Email.RateLimit > 0 {
		config.RateLimit = svcCtx.Config.Email.RateLimit
	}
	if svcCtx.Config.Email.RateBurst > 0 {
		config.RateBurst = svcCtx.Config.Email.RateBurst
	}
	if svcCtx.Config.Email.BatchSize > 0 {
		config.BatchSize = svcCtx.Config.Email.BatchSize
	}

	// Create and start dispatcher
	dispatcher := email.NewDispatcher(sequenceService, svcCtx.DB, config)
	dispatcher.Start()

	return &EmailWorker{
		dispatcher: dispatcher,
	}
}

// Stop gracefully shuts down the email worker
func (w *EmailWorker) Stop() {
	if w.dispatcher != nil {
		w.dispatcher.Stop()
	}
}

// Stats returns current email processing statistics
func (w *EmailWorker) Stats() (sent, failed, retried int64, circuitOpen bool) {
	if w.dispatcher != nil {
		return w.dispatcher.Stats()
	}
	return 0, 0, 0, false
}
