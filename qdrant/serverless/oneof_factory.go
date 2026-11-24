package serverless

// Creates DenseVectors with a single unnamed (default) vector, the one point operations
// use when no vector name is given.
func NewDenseVectors(config *DenseVectorConfig) map[string]*DenseVectorConfig {
	return map[string]*DenseVectorConfig{"": config}
}

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

// Creates a *PayloadIndexConfig instance from *FloatIndex.
// A nil params uses the defaults.
func NewPayloadIndexConfigFloat(params *FloatIndex) *PayloadIndexConfig {
	if params == nil {
		params = &FloatIndex{}
	}
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Float{
			Float: params,
		},
	}
}

// Creates a *PayloadIndexConfig instance from *UuidIndex.
// A nil params uses the defaults.
func NewPayloadIndexConfigUUID(params *UuidIndex) *PayloadIndexConfig {
	if params == nil {
		params = &UuidIndex{}
	}
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Uuid{
			Uuid: params,
		},
	}
}

// Creates a *PayloadIndexConfig instance from *DatetimeIndex.
// A nil params uses the defaults.
func NewPayloadIndexConfigDatetime(params *DatetimeIndex) *PayloadIndexConfig {
	if params == nil {
		params = &DatetimeIndex{}
	}
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Datetime{
			Datetime: params,
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

// Creates a *PayloadIndexConfig instance from *GeoIndex.
// A nil params uses the defaults.
func NewPayloadIndexConfigGeo(params *GeoIndex) *PayloadIndexConfig {
	if params == nil {
		params = &GeoIndex{}
	}
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Geo{
			Geo: params,
		},
	}
}

// Creates a *PayloadIndexConfig instance from *BoolIndex.
// A nil params uses the defaults.
func NewPayloadIndexConfigBool(params *BoolIndex) *PayloadIndexConfig {
	if params == nil {
		params = &BoolIndex{}
	}
	return &PayloadIndexConfig{
		Index: &PayloadIndexConfig_Bool{
			Bool: params,
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
