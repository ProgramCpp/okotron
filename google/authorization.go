package google

import "errors"

var (
	ErrAuthorizationPending = errors.New("")
)

type DeviceCode struct {
	DeviceCode       string `json:device_code`
	Expires_in       int    `json:expires_in`
	Interval         int    `json:interval`
	UserCode         string `json:user_code`
	Verification_url string `json:verification_url`
}

func GetDeviceCode() DeviceCode {
	// reply = fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?response_type=code&client_id=%s&redirect_uri=http://localhost:3000/&scope=https://www.googleapis.com/auth/userinfo.profile&state=123&access_type=offline", os.Getenv("GOOGLE_CLIENT_ID"))

	return DeviceCode{}
}


func PollAuthorization() (string, error) {
	return "", nil
}