package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// natsEnvelope wraps events for NATS transport
type natsEnvelope struct {
	Topic     string          `json:"topic"`
	EmittedAt time.Time       `json:"emitted_at"`
	Payload   json.RawMessage `json:"payload"`
}

// NATSConfig holds NATS connection and configuration options
type NATSConfig struct {
	Conn     *nats.Conn
	Prefix   string
	JSStream string
}

// WithNATS wires NATS into the Subject
func WithNATS(cfg NATSConfig) SubjectOption {
	return func(sc *subjectConfig) {
		sc.natsCfg = cfg
	}
}

// decoderFunc defines how to decode a topic payload into a concrete type
type decoderFunc func([]byte) (any, error)

// RegisterJSONType tells the bus how to decode a topic into a concrete Go type T.
// Call this once during wiring for each typed topic you use cross-process.
func RegisterJSONType[T any](s *Subject, topic string) {
	fn := func(b []byte) (any, error) {
		var v T
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		return v, nil
	}
	for {
		old := s.decoders.Load()
		var nm map[string]decoderFunc
		if old == nil {
			nm = map[string]decoderFunc{topic: fn}
		} else {
			nm = make(map[string]decoderFunc, len(*old)+1)
			for k, v := range *old {
				nm[k] = v
			}
			nm[topic] = fn
		}
		if s.decoders.CompareAndSwap(old, &nm) {
			return
		}
	}
}

// setupNATS initializes NATS connection and subscriptions
func (s *Subject) setupNATS() error {
	if s.config.natsCfg.Conn == nil {
		return nil
	}

	s.nc = s.config.natsCfg.Conn
	s.prefix = s.config.natsCfg.Prefix
	s.natsOn = true

	// Setup JetStream if configured
	if s.config.natsCfg.JSStream != "" {
		js, err := s.nc.JetStream()
		if err != nil {
			return fmt.Errorf("failed to get JetStream context: %w", err)
		}

		// Ensure stream exists
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     s.config.natsCfg.JSStream,
			Subjects: []string{s.prefix + ">"},
			Storage:  nats.FileStorage,
			Replicas: 1,
		})
		if err != nil && err != nats.ErrStreamNameAlreadyInUse {
			return fmt.Errorf("failed to create JetStream: %w", err)
		}

		s.js = js
	}

	// Subscribe to all subjects under prefix and feed local bus
	subj := s.prefix + ">"
	_, err := s.nc.Subscribe(subj, func(m *nats.Msg) {
		var env natsEnvelope
		if err := json.Unmarshal(m.Data, &env); err != nil {
			if s.config.logger != nil {
				s.config.logger.Warn("nats decode error", "err", err)
			}
			return
		}

		// Decode to typed value if we have a decoder; else deliver map[string]any
		dec := (*s.decoders.Load())[env.Topic]
		var msg any
		var decodeErr error
		if dec != nil {
			msg, decodeErr = dec(env.Payload)
		} else {
			var m map[string]any
			decodeErr = json.Unmarshal(env.Payload, &m)
			msg = m
		}
		if decodeErr != nil {
			if s.config.logger != nil {
				s.config.logger.Warn("nats payload decode error", "topic", env.Topic, "err", decodeErr)
			}
			return
		}

		// Inject into the local loop, marking as from NATS to avoid echo
		select {
		case s.events <- event{topic: env.Topic, message: msg, fromNATS: true}:
		default:
			// fall back to blocking to avoid drop
			s.events <- event{topic: env.Topic, message: msg, fromNATS: true}
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to NATS: %w", err)
	}

	return nil
}

// publishToNATS publishes an event to NATS if enabled and not from NATS
func (s *Subject) publishToNATS(evt event) {
	if !s.natsOn || evt.fromNATS {
		return
	}

	b, err := json.Marshal(evt.message)
	if err != nil {
		if s.config.logger != nil {
			s.config.logger.Warn("nats marshal error", "topic", evt.topic, "err", err)
		}
		return
	}

	env := natsEnvelope{
		Topic:     evt.topic,
		EmittedAt: time.Now(),
		Payload:   b,
	}

	data, err := json.Marshal(env)
	if err != nil {
		if s.config.logger != nil {
			s.config.logger.Warn("nats envelope marshal error", "topic", evt.topic, "err", err)
		}
		return
	}

	nsub := s.prefix + evt.topic
	if s.js != nil {
		if _, err := s.js.Publish(nsub, data); err != nil {
			if s.config.logger != nil {
				s.config.logger.Warn("nats js publish error", "subject", nsub, "err", err)
			}
		}
	} else {
		if err := s.nc.Publish(nsub, data); err != nil {
			if s.config.logger != nil {
				s.config.logger.Warn("nats publish error", "subject", nsub, "err", err)
			}
		}
	}
}
