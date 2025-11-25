package acc

type UserInfo struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}
type OAuthResponse struct {
	AccessToken  string `json:"access_token"`            // Bắt buộc
	TokenType    string `json:"token_type"`              // Bắt buộc, thường là "Bearer"
	ExpiresIn    int    `json:"expires_in"`              // Bắt buộc (giây)
	RefreshToken string `json:"refresh_token,omitempty"` // Tùy chọn
	IDToken      string `json:"id_token,omitempty"`      // Tùy chọn (nếu dùng OIDC)
	Scope        string `json:"scope,omitempty"`         // Tùy chọn
}