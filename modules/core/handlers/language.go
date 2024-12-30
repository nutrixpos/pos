package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/gorilla/mux"
)

func GetLanguage(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		lang_code := params["code"]
		pwd, err := os.Getwd()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var foundLanguage map[string]interface{}

		languagesDir := pwd + "/modules/core/languages/"
		files, err := os.ReadDir(languagesDir)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			var languageData map[string]interface{}
			if file.IsDir() {
				continue
			}

			filePath := languagesDir + file.Name()
			jsonFile, err := os.Open(filePath)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			defer jsonFile.Close()

			byteValue, _ := io.ReadAll(jsonFile)
			if err := json.Unmarshal(byteValue, &languageData); err != nil {
				logger.Error(err.Error())
				continue
			}

			if languageData["code"] == lang_code {
				foundLanguage = languageData
			}
		}

		response := JSONApiOkResponse{
			Data: foundLanguage,
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
