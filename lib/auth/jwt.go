package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/k-mistele/ccdc-scoreserver/lib/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)
// AUTH COOKIE LIFESPAN
const authCookieLifespan = 6 * time.Hour

// RANDOMLY GENERATE A SECRET TOKEN
var secretKey = utils.GenerateSecureToken(256)

// DEFINE A TEAM TYPE
type Team string
const (
	Red 	Team = "Red"
	Blue 	Team = "Blue"
	Black 	Team = "Black"
)

// DEFINE CUSTOM JWT CLAIMS
type jwtCustomClaims struct {
	Username 	string		`json:"username"`
	Admin 		bool 		`json:"admin"`
	Team		string 		`json:"team"`
	UUID		string 		`json:"uuid"`
	jwt.StandardClaims
}


// CREATE A FUNCTION TO ISSUE A TOKEN
func NewJSONWebToken(username string, team Team, admin bool, uuid string ) (string, error) {

	// CREATE A TOKEN
	var token *jwt.Token
	var claims jwtCustomClaims
	token = jwt.New(jwt.SigningMethodHS256)

	// SET TOKEN CLAIMS
	claims = jwtCustomClaims{
		username,
		admin,
		string(team),
		uuid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 6).Unix(),
		},
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// SET THE AUTH COOKIE AND MAKE IT LAST FOR HOURS
func SetAuthCookie(c *echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = token
	cookie.Expires = time.Now().Add(authCookieLifespan)
	(*c).SetCookie(cookie)
}

// UNSET THE AUTH COOKIE BY SETTING AN EMPTY, EXPIRED VERSION
func UnsetAuthCookie(c *echo.Context){
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = ""
	cookie.MaxAge = -1
	(*c).SetCookie(cookie)
}