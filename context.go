// Copyright (C) 2017 ScyllaDB

package log

// ctxt is a context key type.
type ctxt byte

// ctxt enumeration.
const (
	ctxTraceID ctxt = iota
)
