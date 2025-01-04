package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/services"
	"github.com/gorilla/mux"
)

// GetLanguage  is an http handler that receives a lang code like "en" or "ar" and reads the
// related language pack json file and return it as response.
func GetLanguage(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		lang_code := params["code"]

		lang_svc := services.LanguageService{
			Logger: logger,
			Config: config,
		}

		lang, err := lang_svc.GetLanguage(lang_code)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONApiOkResponse{
			Data: lang,
		}

		jsonLanguage, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal order settings response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonLanguage)

	}
}

// GetAvailableLanguages reads available language files in the /modules/core/languages folder
// it reads the language code and name from the file and returns json response of available installed langs
func GetAvailableLanguages(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		availableLanguages := []struct {
			Language string `json:"language"`
			Code     string `json:"code"`
		}{}
		pwd, err := os.Getwd()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pwd = pwd + "/modules/core/languages"
		files, err := os.ReadDir(pwd)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filePath := pwd + "/" + file.Name()
			jsonFile, err := os.Open(filePath)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer jsonFile.Close()

			byteValue, _ := io.ReadAll(jsonFile)
			var languageFile struct {
				Language string `json:"language"`
				Code     string `json:"code"`
			}
			json.Unmarshal(byteValue, &languageFile)
			availableLanguages = append(availableLanguages, languageFile)
		}

		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Data []struct {
				Language string `json:"language"`
				Code     string `json:"code"`
			} `json:"data"`
		}{
			Data: availableLanguages,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	}
}
