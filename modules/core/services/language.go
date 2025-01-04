package services

import (
	"encoding/json"
	"io"
	"os"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
)

type LanguageService struct {
	Config   config.Config
	Settings models.Settings
	Logger   logger.ILogger
}

// GetLanguage checks if a language file contains the specified code,
// if so it returns the file
func (ls *LanguageService) GetLanguage(lang_code string) (foundLanguage models.LanguageData, err error) {

	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	languagesDir := pwd + "/modules/core/languages/"
	files, err := os.ReadDir(languagesDir)
	if err != nil {
		return
	}

	for _, file := range files {
		var languageData models.LanguageData
		if file.IsDir() {
			continue
		}

		filePath := languagesDir + file.Name()
		jsonFile, err := os.Open(filePath)
		if err != nil {
			ls.Logger.Error(err.Error())
			continue
		}
		defer jsonFile.Close()

		byteValue, _ := io.ReadAll(jsonFile)
		if err := json.Unmarshal(byteValue, &languageData); err != nil {
			ls.Logger.Error(err.Error())
			continue
		}

		if languageData.Code == lang_code {
			foundLanguage = languageData
		}
	}

	return foundLanguage, nil
}
