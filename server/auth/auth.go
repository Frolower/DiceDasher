package auth

import "dicedasher/storage"

func Auth(access_token string) string {
	id := storage.AccessTokens[access_token]
	return id 
}