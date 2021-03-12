package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/k-mistele/ccdc-scoreserver/lib/messages"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

// ERROR TYPES
var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid")
	ErrMissingToken = errors.New("token is missing")
	ErrForbidden = errors.New("token is for the wrong team")
)


// MIDDLEWARE FOR BLACK TEAM REQUIRED
func BlackTeamRequired(next echo.HandlerFunc) echo.HandlerFunc {

	// RETURN A CALLBACK
	return func(c echo.Context) error {


		// GRAB THE TOKEN FROM THE REQUEST
		jwtCookie, err := getTokenInRequest(&c)
		if err != nil {
			messages.Set(c, messages.Error, "You must log in first!")
			return c.Redirect(http.StatusFound, "/login")
		}

		// PARSE THE TOKEN INTO A *jwt.Token
		token, err := jwt.Parse(jwtCookie, func(token *jwt.Token) (interface{} , error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			log.Error(err)
			UnsetAuthCookie(&c)
			return c.Redirect(http.StatusFound, "/login")
		}
		if !token.Valid {
			log.Debug("Invalid token!")
			UnsetAuthCookie(&c)
			return c.Redirect(http.StatusFound, "/login")
		}

		// MAKE SURE THE TOKEN BELONGS TO THE RIGHT TEAM.
		claims := token.Claims.(jwt.MapClaims)
		log.Debugf("Claims: %v", claims)
		log.Debug(reflect.TypeOf(claims))

		// REQUIRE BLACK TEAM OR ADMIN
		if claims["admin"].(bool) != true  && claims["team"].(string) != string(Black){
			log.Debug("Invalid permissions")
			messages.Set(c, messages.Error, "You don't have permission to do that!")
			return c.Redirect(http.StatusFound, "/")
		}
		return next(c)
	}
}

// RE-IMPLEMENT SOME OF ECHO'S LABSTACK STUFF
func getTokenInRequest (c *echo.Context) (string, error) {
	cookie, err := (*c).Cookie(AuthCookieName)
	if err != nil {
		return "", ErrMissingToken
	}
	return cookie.Value, nil
}
