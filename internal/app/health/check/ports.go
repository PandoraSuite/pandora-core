package check

import "github.com/MAD-py/pandora-core/internal/domain/errors"

type Database interface {
	Ping() errors.Error
	Latency() (int64, errors.Error)
}
