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
	IFramesOnlyTag:           IFramesOnlyParser{},
	ProgramDateTimeTag:       ProgramDateTimeParser{},
	KeyTag:                   KeyParser{},
	MapTag:                   MapParser{},
	DateRangeTag:             DateRangeParser{},
	ExtInfTag:                ExtInfParser{},
	DiscontinuityTag:         DiscontinuityParser{},
	GapTag:                   GapParser{},
	StreamInfTag:             StreamInfParser{},
	MediaTag:                 MediaParser{},
	IFrameStreamInfTag:       IFrameStreamInfParser{},
	SessionKeyTag:            SessionKeyParser{},
	IndependentSegmentsTag:   IndependentSegmentsParser{},
	VariableDefineTag:        VariableDefineParser{},
	USPTimestampMapTag:       USPTimestampMapParser{},
	EventCueOutTag:           EventCueOutParser{},
	EventCueOutContTag:       EventCueOutContParser{},
	EventCueInTag:            EventCueInParser{},
	CommentLineTag:           CommentParser{},
	EndListTag:               EndListParser{},
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
	IFramesOnlyName:           IFramesOnlyEncoder{},
	ProgramDateTimeName:       ProgramDateTimeEncoder{},
	KeyName:                   KeyEncoder{},
	MapName:                   MapEncoder{},
	DateRangeName:             DateRangeEncoder{},
	ExtInfName:                ExtInfEncoder{},
	DiscontinuityName:         DiscontinuityEncoder{},
	GapName:                   GapEncoder{},
	StreamInfName:             StreamInfEncoder{},
	MediaName:                 MediaEncoder{},
	IFrameStreamInfName:       IFrameStreamInfEncoder{},
	SessionKeyName:            SessionKeyEncoder{},
	IndependentSegmentsName:   IndependentSegmentsEncoder{},
	VariableDefineName:        VariableDefineEncoder{},
	USPTimestampMapName:       USPTimestampMapEncoder{},
	EventCueOutName:           EventCueOutEncoder{},
	EventCueOutContName:       EventCueOutContEncoder{},
	EventCueInName:            EventCueInEncoder{},
	CommentLineName:           CommentEncoder{},
	EndListName:               EndListEncoder{},
}
