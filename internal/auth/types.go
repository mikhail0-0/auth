package auth

type User struct {
	GUID string `json:"guid"`
}

type Session struct {
	GUID        string `json:"guid"`
	RefreshHash string `json:"refreshHash"`
	AccessHash string `json:"accessHash"`
	Exp         int64  `json:"exp"`
}

type LoginData struct {
	GUID string `json:"guid"`
}

type RefreshData struct {
	GUID         string `json:"guid"`
	RefreshToken string `json:"refreshToken"`
}

type AccessPayload struct {
	GUID         string `json:"guid"`
	Exp          int64  `json:"exp"`
}

type AccessResponse struct {
	GUID         string `json:"guid"`
	Exp          int64  `json:"exp"`
	AccessToken  string `json:"accessToken"`
}

type TokensData struct {
	GUID         string `json:"guid"`
	Exp          int64  `json:"exp"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
