package length

import (
	"github.com/elastic/beats/libbeat/common"
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstrstr"
	f "github.com/elastic/beats/libbeat/common/schema/mapstriface"
)

var schema = s.Schema{
	"key":    c.Str("key"),
	"length": f.Int("length"),
}

// Map data to MapStr
func eventsMapping(lengthMap map[string]int) []common.MapStr {
	events := []common.MapStr{}
	for key, length := range lengthMap {
		db := map[string]interface{}{}
		db["key"] = key
		db["length"] = length
		data, _ := schema.Apply(db)
		events = append(events, data)
	}

	return events
}
