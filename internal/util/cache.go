package util

import (
	"github.com/elliotchance/orderedmap/v3"
	"github.com/karlseguin/ccache/v3"
)

var CollectionCache *ccache.Cache[*orderedmap.OrderedMap[int64, CollectionData]]
var MapStringCache *ccache.Cache[map[string]string]

func SetupCollectionCache() {
	CollectionCache = ccache.New(ccache.Configure[*orderedmap.OrderedMap[int64, CollectionData]]())
}

func SetupMapStringCache() {
	MapStringCache = ccache.New(ccache.Configure[map[string]string]())
}
