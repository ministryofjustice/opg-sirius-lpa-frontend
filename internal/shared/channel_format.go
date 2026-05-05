package shared

type ChannelFormat int

const (
	ChannelFormatPaper ChannelFormat = iota
	ChannelFormatOnline
	ChannelFormatEmpty
	ChannelFormatNotRecognised
)

var channelFormatMap = map[string]ChannelFormat{
	"paper":         ChannelFormatPaper,
	"online":        ChannelFormatOnline,
	"":              ChannelFormatEmpty,
	"notRecognised": ChannelFormatNotRecognised,
}

func (c ChannelFormat) String() string {
	return c.Key()
}

func (c ChannelFormat) Translation() string {
	switch c {
	case ChannelFormatPaper:
		return "Paper"
	case ChannelFormatOnline:
		return "Online"
	case ChannelFormatEmpty:
		return "Not specified"
	default:
		return "channel NOT RECOGNISED: " + c.String()
	}
}

func (c ChannelFormat) Key() string {
	switch c {
	case ChannelFormatPaper:
		return "Paper"
	case ChannelFormatOnline:
		return "Online"
	case ChannelFormatEmpty:
		return "Empty"
	case ChannelFormatNotRecognised:
		return "Not Recognised"
	default:
		return ""
	}
}

func ParseChannelFormat(s string) ChannelFormat {
	value, ok := channelFormatMap[s]
	if !ok {
		return ChannelFormatNotRecognised
	}
	return value
}

func ChannelForFormat(s string) string {
	return ParseChannelFormat(s).Translation()
}
