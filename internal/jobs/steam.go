package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Age string

func (r *Requirements) UnmarshalJSON(data []byte) error {
	if string(data) == "[]" {
		*r = Requirements{}
		return nil
	}

	type Alias Requirements
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*r = Requirements(tmp)
	return nil
}

func (a *Age) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*a = Age(str)
		return nil
	}
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*a = Age(strconv.Itoa(num))
		return nil
	}
	return fmt.Errorf("invalid age format")
}

type SteamApp struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}

type SteamAppListResponse struct {
	Applist struct {
		Apps []SteamApp `json:"apps"`
	} `json:"applist"`
}

type Requirements struct {
	Minimum     string `json:"minimum"`
	Recommended string `json:"recommended"`
}

type SteamAppDetails struct {
	Success bool `json:"success"`
	Data    struct {
		Name               string       `json:"name"`
		Background         string       `json:"background_raw"`
		HeaderImage        string       `json:"header_image"`
		Age                Age          `json:"required_age"`
		AboutTheGame       string       `json:"about_the_game"`
		Description        string       `json:"detailed_description"`
		ShortDescription   string       `json:"short_description"`
		IsFree             bool         `json:"is_free"`
		Website            string       `json:"website"`
		Legal              string       `json:"legal_notice"`
		SupportedLanguages string       `json:"supported_languages"`
		PCRequirements     Requirements `json:"pc_requirements"`
		MacRequirements    Requirements `json:"mac_requirements"`
		LinuxRequirements  Requirements `json:"linux_requirements"`
		PriceOverview      struct {
			Currency string `json:"currency"`
			Initial  uint   `json:"initial"`
			Final    uint   `json:"final"`
		} `json:"price_overview"`
		ReleaseDate struct {
			Soon bool   `json:"coming_soon"`
			Date string `json:"date"`
		} `json:"release_date"`
		Developers []string `json:"developers"`
		Publishers []string `json:"publishers"`
		Categories []struct {
			Name string `json:"description"`
		}
		Genres []struct {
			Name string `json:"description"`
		} `json:"genres"`
		Support struct {
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Movies []struct {
			MP4 struct {
				Max string `json:"max"`
			} `json:"mp4"`
		} `json:"movies"`
		Screenshots []struct {
			Path string `json:"path_full"`
		} `json:"screenshots"`
		DLC []int `json:"dlc"`
	} `json:"data"`
}

func FetchSteamAppList() ([]SteamApp, error) {
	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	var appListResponse SteamAppListResponse
	if err := json.NewDecoder(resp.Body).Decode(&appListResponse); err != nil {
		return nil, err
	}

	return appListResponse.Applist.Apps, nil
}

func FetchSteamAppDetails(appID int) (*SteamAppDetails, error) {
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%d&cc=en&l=en", appID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	var appDetailsMap map[string]SteamAppDetails
	err = json.NewDecoder(resp.Body).Decode(&appDetailsMap)
	if err != nil {
		return nil, err
	}

	appDetails, ok := appDetailsMap[fmt.Sprintf("%d", appID)]
	if !ok || !appDetails.Success {
		return nil, fmt.Errorf("no data for app ID %d", appID)
	}

	return &appDetails, nil
}
