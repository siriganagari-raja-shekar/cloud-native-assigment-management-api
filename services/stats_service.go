package services

import "github.com/smira/go-statsd"

type StatsStore struct {
	Client *statsd.Client
}
