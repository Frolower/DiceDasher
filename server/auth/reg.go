package auth

import (
	"dicedasher/st"
	"dicedasher/storage"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func Reg (c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(400, nil)
		return 
	}
	var data map[string]string
	json.Unmarshal(jsonData, &data)

	if (data["username"] == "" || data["email"] == "" || data["password_hash"] == "") { 
		c.JSON(400, nil)
		return 
	}
	// Check password, email, username

	var user st.User = st.User{
		Username: data["username"],
		Password_hash: data["password_hash"],
		Email: data["email"],
	}
	user.GenerateID()

	storage.Users[user.Username] = user

	c.JSON(200, nil)
}