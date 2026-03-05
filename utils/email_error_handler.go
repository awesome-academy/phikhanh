package utils

import (
	"log"
	"sync"
	"time"
)

// EmailError - Struct để represent email sending errors
type EmailError struct {
	ApplicationCode string
	RecipientEmail  string
	Error           error
	Timestamp       time.Time
	Retries         int
}

// EmailErrorHandler - Centralized handler cho email errors
type EmailErrorHandler struct {
	errorChan chan *EmailError
	wg        sync.WaitGroup
}

var (
	globalEmailErrorHandler *EmailErrorHandler
	once                    sync.Once
)

// GetEmailErrorHandler - Get singleton instance
func GetEmailErrorHandler() *EmailErrorHandler {
	once.Do(func() {
		globalEmailErrorHandler = &EmailErrorHandler{
			errorChan: make(chan *EmailError, 100), // Buffer untuk 100 errors
		}
		globalEmailErrorHandler.start()
	})
	return globalEmailErrorHandler
}

// start - Start processing email errors
func (h *EmailErrorHandler) start() {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		for err := range h.errorChan {
			h.handleError(err)
		}
	}()
}

// handleError - Process individual email error
func (h *EmailErrorHandler) handleError(emailErr *EmailError) {
	// Log error
	log.Printf("[Email Error] Application: %s, To: %s, Error: %v, Retries: %d",
		emailErr.ApplicationCode,
		emailErr.RecipientEmail,
		emailErr.Error,
		emailErr.Retries,
	)

	// TODO: Implement database logging
	// h.saveToDatabase(emailErr)

	// TODO: Implement admin alert
	// h.sendAdminAlert(emailErr)
}

// ReportError - Report email error
func (h *EmailErrorHandler) ReportError(appCode, recipientEmail string, err error) {
	select {
	case h.errorChan <- &EmailError{
		ApplicationCode: appCode,
		RecipientEmail:  recipientEmail,
		Error:           err,
		Timestamp:       time.Now(),
		Retries:         0,
	}:
	default:
		// Channel full, log directly
		log.Printf("[Email Error] Channel full, dropping error for app %s: %v", appCode, err)
	}
}

// Close - Gracefully shutdown
func (h *EmailErrorHandler) Close() {
	close(h.errorChan)
	h.wg.Wait()
}
