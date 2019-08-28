package domains

import(
	"../utils"
)

type Result struct {
	User 		*User
	Country 	*Country
	Site 		*Site
	ApiError	*utils.ApiError
}