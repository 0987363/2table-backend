package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson/primitive"

	//	"strconv"
	"strings"
	"time"
	//
)

const (
	TokenVersion            = "1"                 // 账号密码
	TokenVersionApplication = "2"                 // 应用id
	tokenExpiryDuration     = time.Hour * 24 * 30 // Valid for 30 days
	tokenNonceLength        = 8
	tokenNonceChars         = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// Token struct represents the secret token
type Token struct {
	//	Base
	Version string    `json:"version" bson:"version"`
	UserID  string    `json:"id" bson:"user_id"`
	Nonce   string    `json:"nonce" bson:"nonce"`
	Expiry  time.Time `json:"expiry" bson:"expiry"`
	Hmac    string    `json:"hmac" bson:"hmac"`
}

// NewToken creates a new Token object from a User object and a secret. The user
// object must have ID field filled.
func NewToken(userID primitive.ObjectID, secret string) *Token {
	token := &Token{
		Version: TokenVersion,
		UserID:  userID.Hex(),
		Nonce:   generateNonce(tokenNonceLength),
		Expiry:  time.Now().Add(tokenExpiryDuration),
	}

	token.Hmac = base64.URLEncoding.EncodeToString(generateHmac(token, secret))
	return token
}

// ParseToken parses the authencation string and returns a pointer to the result
// Token object. If the token is invalid in any way, the pointer will be nil and
// an error will be returned
func ParseToken(token string) (*Token, error) {
	parts := strings.Split(token, ";")
	if len(parts) != 5 {
		return nil, fmt.Errorf("token has %d parts, needs 5", len(parts))
	}

	var err error
	t := &Token{}

	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("token part:%s does not have '=' separator", p)
		}

		switch kv[0] {
		case "v":
			t.Version = kv[1]
		case "id":
			if !IsObjectIDHex(kv[1]) {
				return nil, fmt.Errorf("token user id:%s is not valid", kv[1])
			}
			t.UserID = kv[1]
		case "n":
			t.Nonce = kv[1]
		case "expiry":
			if t.Expiry, err = time.Parse(time.RFC3339, kv[1]); err != nil {
				return nil, fmt.Errorf("token expiry:%s is not valid", kv[1])
			}
		case "hmac":
			t.Hmac = kv[1]
		default:
			return nil, fmt.Errorf("unknown token part: %s", p)
		}
	}

	return t, nil
}

// Convert the token into a URL safe string so it can be used in http header.
func (t *Token) String() string {
	return fmt.Sprintf(
		"v=%s;id=%s;n=%s;expiry=%s;hmac=%s",
		t.Version,
		t.UserID,
		t.Nonce,
		t.Expiry.Format(time.RFC3339),
		t.Hmac,
	)
}

// Validate returns whether the token is an valid token or not.
func (t *Token) Validate(secret string) bool {
	if t.Version != TokenVersion {
		return false
	}

	if !t.Expiry.After(time.Now()) {
		return false
	}

	expectedHmac := generateHmac(t, secret)
	givenHmac, err := base64.URLEncoding.DecodeString(t.Hmac)

	if err != nil {
		return false
	}

	return hmac.Equal(givenHmac, expectedHmac)
}

/*
func (t *Token) CheckUpdated (ctx context.Context, db *qmgo.Client, tokenStr string) bool {
	user := FindUserByID(ctx, db, ObjectIDHex(t.UserID))
	if user == nil {
//		fmt.Println("can't find user:", t.UserID)
		return false
	}

	if user.Token != tokenStr {
//		fmt.Println("token is different:", user.Token, " with ", tokenStr)
		return false
	}

	return true
}
*/

func generateNonce(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	nonce := make([]byte, length)

	for i := 0; i < length; i++ {
		nonce[i] = tokenNonceChars[rand.Intn(len(tokenNonceChars))]
	}

	return string(nonce)
}

func generateHmac(token *Token, secret string) []byte {
	message := fmt.Sprintf(
		"v=%s;id=%s;n=%s;expiry=%s",
		token.Version,
		token.UserID,
		token.Nonce,
		token.Expiry.Format(time.RFC3339),
	)

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(message))
	return hash.Sum(nil)
}
