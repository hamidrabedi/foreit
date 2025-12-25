package models

import (
	"context"
	"sync"
)

// Signal represents a Django-like signal
type Signal struct {
	receivers []SignalReceiver
	mu        sync.RWMutex
}

// SignalReceiver is a function that handles a signal
type SignalReceiver func(ctx context.Context, sender Model, signalType SignalType) error

// SignalType represents the type of signal
type SignalType string

const (
	SignalPreSave    SignalType = "pre_save"
	SignalPostSave   SignalType = "post_save"
	SignalPreDelete  SignalType = "pre_delete"
	SignalPostDelete SignalType = "post_delete"
	SignalPreCreate  SignalType = "pre_create"
	SignalPostCreate SignalType = "post_create"
	SignalPreUpdate  SignalType = "pre_update"
	SignalPostUpdate SignalType = "post_update"
)

// Global signals
var (
	PreSave    = NewSignal()
	PostSave   = NewSignal()
	PreDelete  = NewSignal()
	PostDelete = NewSignal()
	PreCreate  = NewSignal()
	PostCreate = NewSignal()
	PreUpdate  = NewSignal()
	PostUpdate = NewSignal()
)

// NewSignal creates a new signal
func NewSignal() *Signal {
	return &Signal{
		receivers: make([]SignalReceiver, 0),
	}
}

// Connect connects a receiver to the signal
func (s *Signal) Connect(receiver SignalReceiver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.receivers = append(s.receivers, receiver)
}

// Disconnect disconnects a receiver from the signal
func (s *Signal) Disconnect(receiver SignalReceiver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, rec := range s.receivers {
		if &rec == &receiver {
			s.receivers = append(s.receivers[:i], s.receivers[i+1:]...)
			break
		}
	}
}

// Send sends the signal to all receivers
func (s *Signal) Send(ctx context.Context, sender Model, signalType SignalType) error {
	s.mu.RLock()
	receivers := make([]SignalReceiver, len(s.receivers))
	copy(receivers, s.receivers)
	s.mu.RUnlock()
	
	for _, receiver := range receivers {
		if err := receiver(ctx, sender, signalType); err != nil {
			return err
		}
	}
	return nil
}

// ReceiverCount returns the number of connected receivers
func (s *Signal) ReceiverCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.receivers)
}

