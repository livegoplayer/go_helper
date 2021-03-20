package mapstructure

func ToStruct(m map[string]interface{}, s interface{}, tagNames ...string) {
	tagName := "json"
	if len(tagNames) > 0 {
		tagName = tagNames[0]
	}

	config := &DecoderConfig{
		Metadata:         nil,
		Result:           s,
		TagName:          tagName,
		WeaklyTypedInput: true,
	}
	decoder, _ := NewDecoder(config)
	decoder.Decode(m)
}

func ToStructV2(m interface{}, s interface{}, tagNames ...string) {
	tagName := "json"
	if len(tagNames) > 0 {
		tagName = tagNames[0]
	}

	config := &DecoderConfig{
		Metadata:         nil,
		Result:           s,
		TagName:          tagName,
		WeaklyTypedInput: true,
	}
	decoder, _ := NewDecoder(config)
	decoder.Decode(m)
}
