package compression

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type (
	CompressionType  string
	CompressionTypes []CompressionType
)

const (
	NoneCompression   CompressionType = "none"
	GzipCompression                   = "gzip"
	SnappyCompression                 = "snappy"
	Lz4Compression                    = "lz4"
	ZstdCompression                   = "zstd"
)

func Parse(compressionTypes CompressionTypes) kgo.Opt {
	var compressions []kgo.CompressionCodec

	for _, compressionType := range compressionTypes {
		switch compressionType {
		case NoneCompression:
			compressions = append(compressions, kgo.NoCompression())
		case SnappyCompression:
			compressions = append(compressions, kgo.SnappyCompression())
		case ZstdCompression:
			compressions = append(compressions, kgo.ZstdCompression())
		case GzipCompression:
			compressions = append(compressions, kgo.GzipCompression())
		case Lz4Compression:
			compressions = append(compressions, kgo.Lz4Compression())
		default:
			continue
		}
	}

	if len(compressions) == 0 {
		compressions = []kgo.CompressionCodec{kgo.SnappyCompression(), kgo.NoCompression()}
	}

	return kgo.ProducerBatchCompression(compressions...)
}
