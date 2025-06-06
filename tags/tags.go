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

// Parse tag from string to *Node in *Playlist
type TagParser interface {
	Parse(tag string, playlist *pl.Playlist) error
}

var Parsers = map[string]TagParser{
	basic.M3u8IdentifierTag:          basic.M3u8IdentifierParser{},
	basic.VersionTag:                 basic.VersionParser{},
	media.TargetDurationTag:          media.TargetDurationParser{},
	media.MediaSequenceTag:           media.MediaSequenceParser{},
	media.ProgramDateTimeTag:         media.ProgramDateTimeParser{},
	media.DateRangeTag:               media.DateRangeParser{},
	media.ExtInfTag:                  media.ExtInfParser{},
	media.DiscontinuityTag:           media.DiscontinuityParser{},
	multivariant.StreamInfTag:        multivariant.StreamInfParser{},
	exclusive.IndependentSegmentsTag: exclusive.IndependentSegmentsParser{},
	others.USPTimestampMapTag:        others.USPTimestampMapParser{},
	others.EventCueOutTag:            others.EventCueOutParser{},
	others.EventCueInTag:             others.EventCueInParser{},
	others.CommentLineTag:            others.CommentParser{},
}

// Parse tag *Node in *Playlist to string
type PlaylistEncoder interface {
	Encode(node *internal.Node, builder *strings.Builder) error
}

var Encoders = map[string]PlaylistEncoder{
	"M3u8Identifier":      basic.M3u8IdentifierEncoder{},
	"Version":             basic.VersionEncoder{},
	"TargetDuration":      media.TargetDurationEncoder{},
	"MediaSequence":       media.MediaSequenceEncoder{},
	"ProgramDateTime":     media.ProgramDateTimeEncoder{},
	"DateRange":           media.DateRangeEncoder{},
	"ExtInf":              media.ExtInfEncoder{},
	"Discontinuity":       media.DiscontinuityEncoder{},
	"StreamInf":           multivariant.StreamInfEncoder{},
	"IndependentSegments": exclusive.IndependentSegmentsEncoder{},
	"UspTimestampMap":     others.USPTimestampMapEncoder{},
	"CueOut":              others.EventCueOutEncoder{},
	"CueIn":               others.EventCueInEncoder{},
	"Comment":             others.CommentEncoder{},
}
