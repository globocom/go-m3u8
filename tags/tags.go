package tags

import (
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/globocom/go-m3u8/tags/basic"
	"github.com/globocom/go-m3u8/tags/exclusive"
	"github.com/globocom/go-m3u8/tags/media"
	"github.com/globocom/go-m3u8/tags/multivariant"
	"github.com/globocom/go-m3u8/tags/others"
)

// Parse string to *Playlist.
type TagParser interface {
	Parse(tag string, playlist *pl.Playlist) error
}

var Parsers = map[string]TagParser{
	basic.M3u8IdentifierTag:          basic.M3u8IdentifierParser{},
	basic.VersionTag:                 basic.VersionParser{},
	media.TargetDurationTag:          media.TargetDurationParser{},
	media.MediaSequenceTag:           media.MediaSequenceParser{},
	media.DiscontinuitySequenceTag:   media.DiscontinuitySequenceParser{},
	media.ProgramDateTimeTag:         media.ProgramDateTimeParser{},
	media.KeyTag:                     media.ExtKeyParser{},
	media.DateRangeTag:               media.DateRangeParser{},
	media.ExtInfTag:                  media.ExtInfParser{},
	media.DiscontinuityTag:           media.DiscontinuityParser{},
	multivariant.StreamInfTag:        multivariant.StreamInfParser{},
	exclusive.IndependentSegmentsTag: exclusive.IndependentSegmentsParser{},
	exclusive.VariableDefineTag:      exclusive.VariableDefineParser{},
	others.USPTimestampMapTag:        others.USPTimestampMapParser{},
	others.EventCueOutTag:            others.EventCueOutParser{},
	others.EventCueInTag:             others.EventCueInParser{},
	others.CommentLineTag:            others.CommentParser{},
}

// Parse *Playlist to string.
type PlaylistEncoder interface {
	Encode(node *internal.Node, builder *strings.Builder) error
}

var Encoders = map[string]PlaylistEncoder{
	basic.M3u8IdentifierName:          basic.M3u8IdentifierEncoder{},
	basic.VersionName:                 basic.VersionEncoder{},
	media.TargetDurationName:          media.TargetDurationEncoder{},
	media.MediaSequenceName:           media.MediaSequenceEncoder{},
	media.DiscontinuitySequenceName:   media.DiscontinuitySequenceEncoder{},
	media.ProgramDateTimeName:         media.ProgramDateTimeEncoder{},
	media.ExtKeyName:                  media.ExtKeyEncoder{},
	media.DateRangeName:               media.DateRangeEncoder{},
	media.ExtInfName:                  media.ExtInfEncoder{},
	media.DiscontinuityName:           media.DiscontinuityEncoder{},
	multivariant.StreamInfName:        multivariant.StreamInfEncoder{},
	exclusive.IndependentSegmentsName: exclusive.IndependentSegmentsEncoder{},
	exclusive.VariableDefineName:      exclusive.VariableDefineEncoder{},
	others.USPTimestampMapName:        others.USPTimestampMapEncoder{},
	others.EventCueOutName:            others.EventCueOutEncoder{},
	others.EventCueInName:             others.EventCueInEncoder{},
	others.CommentLineName:            others.CommentEncoder{},
}
