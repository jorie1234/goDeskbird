package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	login := os.Getenv("LOGIN")
	pwd := os.Getenv("PWD")

	if len(login) == 0 || len(pwd) == 0 {
		fmt.Println("No login or pwd")
		return
	}
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImY5N2U3ZWVlY2YwMWM4MDhiZjRhYjkzOTczNDBiZmIyOTgyZTg0NzUiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vZGVza2JpcmQtYmJlNzIiLCJhdWQiOiJkZXNrYmlyZC1iYmU3MiIsImF1dGhfdGltZSI6MTY4Njg1MzQ2MCwidXNlcl9pZCI6InlETnRVbUp3eWZicXNkZXlOVHE3YVlYWk9OcjEiLCJzdWIiOiJ5RE50VW1Kd3lmYnFzZGV5TlRxN2FZWFpPTnIxIiwiaWF0IjoxNjg4MjI1MjY0LCJleHAiOjE2ODgyMjg4NjQsImVtYWlsIjoiam9uYXMucmllZGVsQHZpZXIuYWkiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJtaWNyb3NvZnQuY29tIjpbImNiZWQ1OGRlLTlmYjUtNDE3MC04OGI0LTFjNDdmYTc2NjQwNiJdLCJlbWFpbCI6WyJqb25hcy5yaWVkZWxAdmllci5haSJdfSwic2lnbl9pbl9wcm92aWRlciI6Im1pY3Jvc29mdC5jb20ifX0.KIykwvBsHJpI5TBWLlX-iz4qjj73THdSHqnKiFLmFZ1z3xyKzCKHeodxSFLk_pIqMHvj6fF1xk-9n_yGyqe6xbvpCtvQtemV6DZoefUC6MiczI3tccBzoicio0iOHCfFQCFUOd_NpXJVJD2tyEarELXKzpTlrn_2lN7GKf-4kVe-abN7W4-IoaXafhNthJUIuY-B9yMqaJuzFlfoEb1cCHGEItLxdTL0gkEwqc1hHi2Y2u7ReOt71BoZhs4xWqFpUf8zUy1KXcoDSVn9SzFkzI-T6qmxq39omHagvFvuogTQf_mYo1A-9QpbT67qUfBnio4jkaq_gk76u8CWEdBzLA"
	err := GetUser(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = LoadLogin(token)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func LoadLogin(token string) error {
	cl := *http.DefaultClient
	cl.Timeout = 30 * time.Second
	cl.Jar = requests.NewCookieJar()

	for i := 1; i < 60; i++ {
		t := time.Now().AddDate(0, 0, -i)

		tFormat := t.Format("2006-01-02")
		tFileName := t.Format("2006-01-02.json")

		//check if file exists
		if _, err := os.Stat(tFileName); !os.IsNotExist(err) {
			continue
		}
		u := "https://app.deskbird.com/api/v1.1/officePlanning?companyId=938&day=" + tFormat + "&allUsers=trun"
		var db Deskbird
		err := requests.
			URL(u).
			Client(&cl).
			Header("authorization", "Bearer "+token).
			Handle(requests.ToJSON(&db)).
			Transport(requests.Record(nil, "")).
			Fetch(context.Background())
		if err != nil {
			return err
		}

		fmt.Printf("%+#v", db)
		//save to file
		f, err := os.Create(tFileName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer f.Close()

		b, err := json.Marshal(db)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = f.Write(b)
		if err != nil {
			fmt.Println(err)
			return err
		}

	}
	return nil

}

func GetUser(token string) error {
	cl := *http.DefaultClient
	cl.Timeout = 30 * time.Second
	cl.Jar = requests.NewCookieJar()

	tFileName := time.Now().Format("users2006-01-02.json")

	u := "https://app.deskbird.com/api/v1.1/businesscompany/users?companyId=938"
	var db UserData
	err := requests.
		URL(u).
		Client(&cl).
		Header("authorization", "Bearer "+token).
		Handle(requests.ToJSON(&db)).
		Transport(requests.Record(nil, "")).
		Fetch(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("%+#v", db)
	//save to file
	f, err := os.Create(tFileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	b, err := json.Marshal(db)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

type Deskbird struct {
	Success bool       `json:"success"`
	Data    DataGLobal `json:"data"`
}
type SelectedOption struct {
	OfficeID   string `json:"officeId"`
	OfficeName string `json:"officeName"`
	OptionID   string `json:"optionId"`
}

type Data struct {
	ID              string         `json:"id"`
	FirstName       string         `json:"firstName"`
	LastName        string         `json:"lastName"`
	Email           string         `json:"email"`
	AvatarColor     string         `json:"avatarColor"`
	ProfileImage    string         `json:"profileImage"`
	SelectedOption  SelectedOption `json:"selectedOption,omitempty"`
	IsUserFavourite bool           `json:"isUserFavourite"`
}
type Starred struct {
	Data  []Data `json:"data"`
	Total int    `json:"total"`
}
type All struct {
	Data  []Data `json:"data"`
	Total int    `json:"total"`
}
type DataGLobal struct {
	Starred Starred `json:"starred"`
	All     All     `json:"all"`
}

type UserData struct {
	Success bool      `json:"success"`
	Results []Results `json:"results"`
}
type UserSettings struct {
	EnableCalendarInvites              bool `json:"enableCalendarInvites"`
	EnableCheckInReminderPN            bool `json:"enableCheckInReminderPN"`
	EnableScheduleReminderPN           bool `json:"enableScheduleReminderPN"`
	EnableCheckInReminderEmail         bool `json:"enableCheckInReminderEmail"`
	EnableScheduleReminderEmail        bool `json:"enableScheduleReminderEmail"`
	EnableAutoCancellationNoticePN     bool `json:"enableAutoCancellationNoticePN"`
	EnableOfficePlanningForOthersPN    bool `json:"enableOfficePlanningForOthersPN"`
	EnableAutoCancellationNoticeEmail  bool `json:"enableAutoCancellationNoticeEmail"`
	EnableOfficePlanningForOthersEmail bool `json:"enableOfficePlanningForOthersEmail"`
}
type ExternalUserData struct {
	ID                string `json:"id"`
	Provider          string `json:"provider"`
	SyncExternalImage bool   `json:"syncExternalImage"`
}
type Results struct {
	ID                    string           `json:"id"`
	FirstName             string           `json:"firstName"`
	LastName              string           `json:"lastName"`
	Email                 string           `json:"email"`
	FcmTokens             []interface{}    `json:"fcmTokens"`
	AvatarColor           string           `json:"avatarColor"`
	UserGroupIdsFirebase  []interface{}    `json:"userGroupIdsFirebase"`
	ProfileImage          string           `json:"profileImage"`
	CompanyID             string           `json:"companyId"`
	HealthCheckRequired   bool             `json:"healthCheckRequired"`
	ExpirationDate        time.Time        `json:"expirationDate"`
	Status                string           `json:"status"`
	Role                  string           `json:"role"`
	Signup                bool             `json:"signup"`
	DemoUser              bool             `json:"demoUser"`
	UserSettings          UserSettings     `json:"userSettings"`
	ExternalUserData      ExternalUserData `json:"externalUserData"`
	PrimaryOfficeID       string           `json:"primaryOfficeId"`
	FavoriteDesks         []interface{}    `json:"favoriteDesks"`
	Favourites            []int            `json:"favourites"`
	RoleLastChangedBy     string           `json:"roleLastChangedBy"`
	CreatedAt             time.Time        `json:"createdAt"`
	UpdatedAt             time.Time        `json:"updatedAt"`
	ExcludeFromPlannings  bool             `json:"excludeFromPlannings"`
	Language              string           `json:"language"`
	FirebaseID            string           `json:"firebaseId"`
	InitialDeviceLanguage string           `json:"initialDeviceLanguage"`
	IsUsingSystemLanguage bool             `json:"isUsingSystemLanguage"`
	UUID                  string           `json:"uuid"`
	UserGroupIds          []string         `json:"userGroupIds"`
}
