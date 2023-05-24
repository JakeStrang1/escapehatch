package db

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

/****************************************************************************************
 * registry.go
 *
 * This file is intended to:
 * - Provide a custom struct codec that swaps the "bson" struct tag for a "db" tag
 ****************************************************************************************/

func Registry() *bsoncodec.Registry {
	rb := bson.NewRegistryBuilder()
	rb.RegisterDefaultDecoder(reflect.Struct, newCustomStructCodec())
	rb.RegisterDefaultEncoder(reflect.Struct, newCustomStructCodec())
	return rb.Build()
}

func newCustomStructCodec() *bsoncodec.StructCodec {
	codec, err := bsoncodec.NewStructCodec(CustomStructTagParser)
	if err != nil {
		// This function is called from the codec registration path, so errors can't be propagated. If there's an error
		// constructing the StructCodec, we panic to avoid losing it.
		panic(fmt.Errorf("error creating StructCodec: %v", err))
	}
	return codec
}

// CustomStructTagParser is a copy of the default StructTagParser used by the StructCodec in bsoncodec/default_value_encoders.go
// It is exactly the same as the original except the tag "bson" is replaced with "db"
// Motivation: I thought it would be nice to abstract away as many mongo-related details from the code as I could.
var CustomStructTagParser bsoncodec.StructTagParserFunc = func(sf reflect.StructField) (bsoncodec.StructTags, error) {
	key := strings.ToLower(sf.Name)
	tag, ok := sf.Tag.Lookup("db")
	if !ok && !strings.Contains(string(sf.Tag), ":") && len(sf.Tag) > 0 {
		tag = string(sf.Tag)
	}
	return parseTags(key, tag)
}

func parseTags(key string, tag string) (bsoncodec.StructTags, error) {
	var st bsoncodec.StructTags
	if tag == "-" {
		st.Skip = true
		return st, nil
	}

	for idx, str := range strings.Split(tag, ",") {
		if idx == 0 && str != "" {
			key = str
		}
		switch str {
		case "omitempty":
			st.OmitEmpty = true
		case "minsize":
			st.MinSize = true
		case "truncate":
			st.Truncate = true
		case "inline":
			st.Inline = true
		}
	}

	st.Name = key

	return st, nil
}
