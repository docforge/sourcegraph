package workqueue

import "errors"

// ErrDequeueTransaction occurs when Dequeue is called from inside a transaction.
var ErrDequeueTransaction = errors.New("unexpected transaction")

// ErrDequeueRace occurs when an upload selected for dequeue has been locked by another worker.
var ErrDequeueRace = errors.New("dequeue race")

// ErrNoRecord occurs when a record cannot be selected after it ha been locked.
var ErrNoRecord = errors.New("no record")
