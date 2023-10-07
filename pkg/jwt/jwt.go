package jwt

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

var (
	// Generated with: ssh-keygen -t rsa -f key.pem -m pem
	privatePEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
	publicPEM = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2Tu9OBa3aMbPYe3odDXK
gLQhqc5z6S1H4Z+svq73uVKEDu49k22V42mr7aAw5sQfylBtu2jXIejM2WNGfyqf
9qRq7wcv4GeOfcYv7EjYp/B5b4HKazxlm8YncU5/ZT7kUizr0D09rPUQSnPUXyTj
w6SMZ2UpKN8aB1KAK7Rxk3KD+sPA+7NLEG7r3XDvKrXLyIl//eEVEYiOma4g2ZCn
t8R61aVNw0sjLaxiYsiJ8uau0poCyAreTrIJAPDOXhQWwUsKNyHcaQgtsz/l+8Ke
GWWrj7u6iAGfRH8lzZdtoHzlLuRwIK4D3P9DyjORwQ7kwzo+D/A+0UlhLxHvwhZf
Q/hHHbtjoWglIDTmAIWUX/eDbGLPXkyty2J2R79ekEBG5bd9gWH9k5ev+oMYeCu7
FLifjzJC3Jx53cz4BsUYVmtzGdGYYga2MD/0aSdOcM/A64pSTy3d2+WPKycotTcq
n+KZEK9jI4sogD3GKYsrlP6r3Ja5/25jg5E0H8VmfJtJjmiF4DR6IBVIIUslYkZI
o56TBghUFSvZc939A+jSVUsqimpxBdKKNwiufXh4ZkMd/D3mpypHfTo/aWvB20nU
9dEusb/9JpWNKEyfMeJX+AukjQeJNFBrV8m152DlEKSeGJANgb9cd8Rnc+0+s+aE
GxdO8meq2DOtx1+LQMQgCAkCAwEAAQ==
-----END PUBLIC KEY-----`)
)

var ErrInvalidToken = errors.New("parsed token is invalid")

type Claims struct {
	Role    string
	Subject string
	Pii     []byte
}

type PublicClaims struct {
	Role string `json:"rol,omitempty"`
	PII  []byte `json:"pii,omitempty"`
}

type claimsWrapper struct {
	PublicClaims
	jwt.RegisteredClaims
}

type option func(*jwt.RegisteredClaims)

func WithExpiration(time time.Time) option {
	return func(claims *jwt.RegisteredClaims) {
		claims.ExpiresAt = jwt.NewNumericDate(time)
	}
}

func ForActor(id string) option {
	return func(claims *jwt.RegisteredClaims) {
		claims.Subject = id
	}
}

func createToken(publicClaims PublicClaims, options ...option) *jwt.Token {
	registeredClaims := jwt.RegisteredClaims{
		ID:       ulid.Make().String(),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		/* If the principal processing the claim does not identify
		itself with a value in the "aud" claim when this claim
		is present, then the JWT MUST be rejected.*/
		Audience: jwt.ClaimStrings{"flat"},
	}
	for _, option := range options {
		option(&registeredClaims)
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claimsWrapper{
		PublicClaims:     publicClaims,
		RegisteredClaims: registeredClaims,
	})
}

func sign(token *jwt.Token) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func Create(publicClaims PublicClaims, options ...option) (string, error) {
	return sign(createToken(publicClaims, options...))
}

func CreateAccessCookie(publicClaims PublicClaims, options ...option) (http.Cookie, error) {
	return CreateCookie("access", publicClaims, options...)
}

func CreateCookie(name string, publicClaims PublicClaims, options ...option) (http.Cookie, error) {
	cookie := http.Cookie{
		Name:     name,
		HttpOnly: true,
		Path:     "/",
	}

	// We have to sign it in this method and not use Create
	// As we need the expiration time of the token for the cookie.
	token := createToken(publicClaims, options...)
	signed, err := sign(token)
	if err != nil {
		return cookie, err
	}

	cookie.Value = signed

	if expiration, err := token.Claims.GetExpirationTime(); err == nil {
		cookie.Expires = expiration.Time
	}

	return cookie, nil
}

func Parse(raw string) (Claims, error) {
	var claims Claims

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return claims, err
	}

	token, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			panic("Invalid or unkown jwt signing method")
		}
		return publicKey, nil
	})
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return claims, ErrInvalidToken
	}

	mappedClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		panic("Token claims could not be mapped")
	}

	claims.Subject = mappedClaims["sub"].(string)
	claims.Role = mappedClaims["rol"].(string)

	// byte arrays are serialised to base64 and therfore must be decoded.
	if pii64, exists := mappedClaims["pii"]; exists {
		pii, err := base64.StdEncoding.DecodeString(pii64.(string))
		if err != nil {
			return claims, err
		}

		claims.Pii = pii
	}

	return claims, nil
}
