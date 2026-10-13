package serverless

// Creates a *PayloadIndexConfig instance from *KeywordIndex.
// A nil params uses the defaults.
func NewPayloadIndexConfigKeyword(params *KeywordIndex) *PayloadIndexConfig {
	if params == nil {
		params = &KeywordIndex{}
	}
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Keyword{
			Keyword: params,
		},
	}
}

// Creates a *PayloadIndexConfig instance from *IntegerIndex.
// A nil params uses the defaults.
func NewPayloadIndexConfigInteger(params *IntegerIndex) *PayloadIndexConfig {
	if params == nil {
		params = &IntegerIndex{}
	}
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Integer{
			Integer: params,
		},
	}
}

// Creates a *PayloadIndexConfig instance for a float index.
func NewPayloadIndexConfigFloat() *PayloadIndexConfig {
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Float{
			Float: &FloatIndex{},
		},
	}
}

// Creates a *PayloadIndexConfig instance for a UUID index.
func NewPayloadIndexConfigUUID() *PayloadIndexConfig {
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Uuid{
			Uuid: &UuidIndex{},
		},
	}
}

// Creates a *PayloadIndexConfig instance for a datetime index.
func NewPayloadIndexConfigDatetime() *PayloadIndexConfig {
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Datetime{
			Datetime: &DatetimeIndex{},
		},
	}
}

// Creates a *PayloadIndexConfig instance from *TextIndex.
// A nil params uses the defaults.
func NewPayloadIndexConfigText(params *TextIndex) *PayloadIndexConfig {
	if params == nil {
		params = &TextIndex{}
	}
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Text{
			Text: params,
		},
	}
}

// Creates a *PayloadIndexConfig instance for a geo index.
func NewPayloadIndexConfigGeo() *PayloadIndexConfig {
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Geo{
			Geo: &GeoIndex{},
		},
	}
}

// Creates a *PayloadIndexConfig instance for a bool index.
func NewPayloadIndexConfigBool() *PayloadIndexConfig {
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Bool{
			Bool: &BoolIndex{},
		},
	}
}

// Creates a *StemmingAlgorithm instance for the snowball stemmer, e.g. for "english".
func NewStemmingAlgorithmSnowball(language string) *StemmingAlgorithm {
	return &StemmingAlgorithm{
		StemmingParams: &StemmingAlgorithm_Snowball{
			Snowball: &SnowballParams{
				Language: language,
			},
		},
	}
}

// Creates a *StemmingAlgorithm instance that explicitly disables stemming.
func NewStemmingAlgorithmDisabled() *StemmingAlgorithm {
	return &StemmingAlgorithm{
		StemmingParams: &StemmingAlgorithm_Disabled{
			Disabled: &DisabledStemmer{},
		},
	}
}
