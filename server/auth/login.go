package auth

import (
	"crypto/sha256"
	"dicedasher/storage"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

func createAccessToken(username string, password string) string {
	accessString := username + password + string(time.Now().UnixNano())
	accessToken_bytes := sha256.New()
	accessToken_bytes.Write([]byte(accessString))
	accessToken :=  hex.EncodeToString(accessToken_bytes.Sum(nil))
	return accessToken
}

func Login(c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(400, nil)
		return 
	}
	var data map[string]string
	json.Unmarshal(jsonData, &data)

	if (data["username"] == "" || data["password"] == "") { 
		c.JSON(400, nil)
		return 
	}

	password_hash_bytes := sha256.New()
	password_hash_bytes.Write([]byte(data["password"]))
	password_hash := hex.EncodeToString(password_hash_bytes.Sum(nil))

	username := data["username"]

	// fmt.Println(data["email"], string(password_hash))
	// fmt.Println(storage.Users[data["email"]], storage.Users[data["email"]].Password_hash)
	// fmt.Println(password_hash, storage.Users[email].Password_hash)

	if (password_hash == storage.Users[username].Password_hash) {
		token := createAccessToken(username, password_hash)
		storage.AccessTokens[token] = storage.Users[username].Id
		c.SetCookie("access_token", token, 1000000000000, "/", "localhost", true, true)
		c.JSON(200, nil)
	} else {
		c.JSON(401, nil)
	}

}