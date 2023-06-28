package util

import (
	"github.com/netrixframework/netrix/strategies/fuzzing"
	"github.com/netrixframework/netrix/types"
)

func TLCEventMapper() fuzzing.TLCMapper {
	// TODO: complete this
	defaultMapper := fuzzing.DefaultEventMapper()
	return func(l *types.List[*types.Event]) []fuzzing.TlcEvent {

		return defaultMapper(l)
	}
}
