package models

// LanguageData is used to store the language mentioned details,
// stored in the language json file.
type LanguageData struct {
	ApiVersion  string                 `json:"api_version" bson:"api_version" mapstructure:"api_version"`
	Code        string                 `json:"code" bson:"code" mapstructure:"code"`
	Language    string                 `json:"language" bson:"language" mapstructure:"language"`
	Orientation string                 `json:"orientation" bson:"orientation" mapstructure:"orientation"`
	Pack        map[string]interface{} `json:"pack" bson:"pack" mapstructure:"pack"`
}
