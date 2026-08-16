package time_provider

import "time"

type TimeProvider struct{}

func New() *TimeProvider {
	return &TimeProvider{}
}

func (p *TimeProvider) Now() time.Time {
	return p.Normalize(time.Now())
}

func (p *TimeProvider) Normalize(t time.Time) time.Time {
	return t.UTC().Truncate(time.Microsecond)
}
