package tags

import (
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
)

// Parse string to *Playlist.
type TagParser interface {
	Parse(tag string, playlist *pl.Playlist) error
}

var Parsers = map[string]TagParser{
	M3u8IdentifierTag:        M3u8IdentifierParser{},
	VersionTag:               VersionParser{},
	TargetDurationTag:        TargetDurationParser{},
	MediaSequenceTag:         MediaSequenceParser{},
	DiscontinuitySequenceTag: DiscontinuitySequenceParser{},
	ProgramDateTimeTag:       ProgramDateTimeParser{},
	KeyTag:                   ExtKeyParser{},
	DateRangeTag:             DateRangeParser{},
	ExtInfTag:                ExtInfParser{},
	DiscontinuityTag:         DiscontinuityParser{},
	StreamInfTag:             StreamInfParser{},
	MediaTag:                 MediaParser{},
	IndependentSegmentsTag:   IndependentSegmentsParser{},
	VariableDefineTag:        VariableDefineParser{},
	USPTimestampMapTag:       USPTimestampMapParser{},
	EventCueOutTag:           EventCueOutParser{},
	EventCueInTag:            EventCueInParser{},
	CommentLineTag:           CommentParser{},
}

// Parse *Playlist to string.
type PlaylistEncoder interface {
	Encode(node *internal.Node, builder *strings.Builder) error
}

var Encoders = map[string]PlaylistEncoder{
	M3u8IdentifierName:        M3u8IdentifierEncoder{},
	VersionName:               VersionEncoder{},
	TargetDurationName:        TargetDurationEncoder{},
	MediaSequenceName:         MediaSequenceEncoder{},
	DiscontinuitySequenceName: DiscontinuitySequenceEncoder{},
	ProgramDateTimeName:       ProgramDateTimeEncoder{},
	ExtKeyName:                ExtKeyEncoder{},
	DateRangeName:             DateRangeEncoder{},
	ExtInfName:                ExtInfEncoder{},
	DiscontinuityName:         DiscontinuityEncoder{},
	StreamInfName:             StreamInfEncoder{},
	MediaName:                 MediaEncoder{},
	IndependentSegmentsName:   IndependentSegmentsEncoder{},
	VariableDefineName:        VariableDefineEncoder{},
	USPTimestampMapName:       USPTimestampMapEncoder{},
	EventCueOutName:           EventCueOutEncoder{},
	EventCueInName:            EventCueInEncoder{},
	CommentLineName:           CommentEncoder{},
}
