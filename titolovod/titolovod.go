// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package titolovod

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const searchURLbeginning = "https://www.timvision.it/TIM/10.14.11/PROD/IT/CUBOWEB/ITALY/DETAILS?contentId="
const searchURLending = "&deviceType=ALL&serviceName=ALL&type=ALL"

// Response riporta i dati della richiesta API.
type Response []struct {
	ResultCode string        `json:"resultCode"`
	SystemTime int           `json:"systemTime"`
	Errors     []interface{} `json:"errors"`
	ResultObj  struct {
		ID         string `json:"id"`
		Containers []struct {
			Actions []struct {
				Key        string `json:"key"`
				URI        string `json:"uri"`
				TargetType string `json:"targetType"`
			} `json:"actions"`
			Layout   string `json:"layout"`
			Metadata struct {
				ShortDescription      string      `json:"-"`
				ParentalControlLevel  string      `json:"parentalControlLevel"`
				Language              string      `json:"language"`
				ContentID             string      `json:"contentId"`
				ExpirationDate        int         `json:"expirationDate"`
				Title                 string      `json:"title"`
				ContentAnalyticsType  string      `json:"contentAnalyticsType"`
				Type                  string      `json:"type"`
				LongDescription       string      `json:"longDescription"`
				Year                  string      `json:"year"`
				ContentProvider       string      `json:"contentProvider"`
				Duration              int         `json:"duration"`
				Genre                 string      `json:"genre"`
				SupportedDevices      []string    `json:"supportedDevices"`
				CategoriesThemeArea   []string    `json:"categoriesThemeArea"`
				Rating                string      `json:"rating"`
				Badge                 string      `json:"badge"`
				OfferTypeLabel        interface{} `json:"offerTypeLabel"`
				ImageURL              string      `json:"imageUrl"`
				ImageURLOther         string      `json:"imageUrlOther"`
				BgImageURL            string      `json:"bgImageUrl"`
				DisplayExpirationDate string      `json:"displayExpirationDate"`
				Directors             []string    `json:"directors"`
				Actors                []string    `json:"actors"`
				BroadcastDisplayDate  string      `json:"broadcastDisplayDate"`
				ObjectSubtype         string      `json:"objectSubtype"`
				ObjectType            string      `json:"objectType"`
				TitleBrief            string      `json:"titleBrief"`
				VideoType             []string    `json:"videoType"`
			}
		}
	}
}

// Get recupera i dati del VOD
func Get(ctx context.Context, idvideoteca string) (result Response, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	URLPAREMETRIZZATA := searchURLbeginning + idvideoteca + searchURLending

	//fmt.Println(URLPAREMETRIZZATA)

	resp, err := http.Get(URLPAREMETRIZZATA)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	result, err = ElaboraResp(bodyBytes)

	return result, err
}


// ElaboraResp elabora la risposta ricevuta dal server.
func ElaboraResp(bodyBytes []byte) (result Response, err error) {

	// trasforma in stringa
	bodyString := string(bodyBytes)

	// Elimina il carattere di carriage return di windons che
	// ogni tanto si trova nella descrizione dei vod.
	bodyString2 := strings.ReplaceAll(bodyString, "\u003f", "")


	err = json.Unmarshal([]byte(bodyString2), &result)
	if err != nil {
		log.Println(err)
	}

	return result, err
}

//  0x000A, 0x000B, 0x000C, 0x000D, 0x0085, 0x2028, 0x2029
