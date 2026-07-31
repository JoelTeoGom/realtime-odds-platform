package domain

import "time"

type EventStatus string

const (
	EventScheduled EventStatus = "scheduled"
	EventLive      EventStatus = "live"
	EventFinished  EventStatus = "finished"
	EventSuspended EventStatus = "suspended"
)

type MarketStatus string

const (
	MarketOpen      MarketStatus = "open"
	MarketSuspended MarketStatus = "suspended"
	MarketClosed    MarketStatus = "closed"
)

type OddStatus string

const (
	OddActive    OddStatus = "active"
	OddSuspended OddStatus = "suspended"
)

type Event struct {
	ID        string
	Name      string
	Sport     string
	StartTime time.Time
	Status    EventStatus
	Markets   map[string]*Market
}

type Market struct {
	ID      string
	EventID string
	Type    string
	Status  MarketStatus
	Odds    map[string]*Odd
}

type Odd struct {
	ID       string
	MarketID string
	Outcome  string
	Price    float64
	Status   OddStatus
	Version  int64
}

type OddUpdate struct {
	EventID  string
	MarketID string
	OddID    string
	Price    *float64
	Status   *OddStatus
	Version  int64
}
