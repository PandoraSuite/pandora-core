package reservation

import (
	"github.com/MAD-py/pandora-core/internal/app/reservation/commit"
	"github.com/MAD-py/pandora-core/internal/app/reservation/rollback"
)

// ... Commit Use Case ...

type ReservationCommitRepository = commit.ReservationRepository

// ... Rollback Use Case ...

type ReservationRollbackRepository = rollback.ReservationRepository
type EnvironmentAvailableRequestIncrementerRepository = rollback.EnvironmentRepository
