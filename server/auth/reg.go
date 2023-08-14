package auth

import (
	"dicedasher/actions"
	"dicedasher/st/db"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Reg(c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(400, nil)
		return
	}
	var data map[string]string
	json.Unmarshal(jsonData, &data)

	if data["username"] == "" || data["email"] == "" || data["password_hash"] == "" {
		c.JSON(400, nil)
		return
	}
	// Check password, email, username

	var user = db.User{
		Username:      data["username"],
		Password_hash: data["password_hash"],
		Email:         data["email"],
	}
	//user.GenerateID()

	if actions.EstablishConnection() {
		defer actions.CloseConnection(db.UserClient, db.UserContext)

		err = db.CreateUser(c, db.UserClient, user)
		if err != nil {
			fmt.Println(err)
			fmt.Println("can't create a user")
			c.JSON(500, nil)
			return
		}

		c.JSON(200, nil)
	} else {
		c.JSON(500, nil)
	}
}
