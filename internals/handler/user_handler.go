package handler

import (
	"net/http"
	datamodels "rest-api/rest-api/datamodels"
	middleware "rest-api/rest-api/internals/middleware/auth"
	"rest-api/rest-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

var users []datamodels.User

func GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	for _, u := range users {
		if u.UID == id {
			ctx.JSON(
				http.StatusFound,
				gin.H{
					"User": u,
				},
			)
			return
		}
	}
	ctx.JSON(
		http.StatusNotFound,
		gin.H{
			"User": "nil",
		},
	)
}

func GetAllUser(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"users": users,
		},
	)
}

func CreateUser(ctx *gin.Context) {
	user := datamodels.User{}
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
		return
	}

	if user.Email == "" || user.Name == "" || user.Password == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "user name/email/password cannot be empty",
			},
		)
		return
	}
	uid, err := utils.GenerateUID()
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err,
			},
		)
		return
	}

	user.UID = uid
	users = append(users, user)
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"user": user,
		},
	)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	for i, u := range users {
		if u.UID == id {
			if u.IsDeleted {
				users = append(users[:i], users[i+1:]...)
				ctx.JSON(
					http.StatusFound,
					gin.H{
						"User":   u,
						"status": "deleted",
					},
				)
				return
			}
			users[i].IsDeleted = true
			ctx.JSON(
				http.StatusFound,
				gin.H{
					"User":   u,
					"status": "marked deleted",
				},
			)
			return
		}
	}
	ctx.JSON(
		http.StatusNotFound,
		gin.H{
			"User": "nil",
		},
	)
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")

	userUpdateReq := datamodels.UserUpdateReq{}
	if err := ctx.ShouldBindBodyWithJSON(&userUpdateReq); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
		return
	}

	for i, u := range users {
		if u.UID == id {
			users[i].Age = userUpdateReq.Age
			users[i].Name = userUpdateReq.Name
			users[i].Email = userUpdateReq.Email
			users[i].Gender = userUpdateReq.Gender

			ctx.JSON(
				http.StatusOK,
				gin.H{
					"User": users[i],
				},
			)
			return
		}
	}

	ctx.JSON(
		http.StatusNotFound,
		gin.H{
			"User": "nil",
		},
	)

}

func Patch(ctx *gin.Context) {
	id := ctx.Param("id")

	userUpdateReq := datamodels.UserUpdateReq{}
	if err := ctx.ShouldBindBodyWithJSON(&userUpdateReq); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
		return
	}

	for i, u := range users {
		if u.UID == id {
			if userUpdateReq.Age != 0 {
				users[i].Age = userUpdateReq.Age
			}
			if userUpdateReq.Gender != "" {
				users[i].Gender = userUpdateReq.Gender
			}
			if userUpdateReq.Name != "" {
				users[i].Name = userUpdateReq.Name
			}
			if userUpdateReq.Email != "" {
				users[i].Email = userUpdateReq.Email
			}

			ctx.JSON(
				http.StatusOK,
				gin.H{
					"User": users[i],
				},
			)
			return
		}
	}

	ctx.JSON(
		http.StatusNotFound,
		gin.H{
			"User": "nil",
		},
	)

}

func Login(ctx *gin.Context) {
	loginReq := datamodels.LoginReq{}
	if err := ctx.ShouldBindBodyWithJSON(&loginReq); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
		return
	}
	if loginReq.Pwd == "" || loginReq.UserId == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "password and userid are required",
			},
		)
		return
	}

	for _, u := range users {
		if u.UID == loginReq.UserId {
			if u.Password == loginReq.Pwd {
				tokens, err := middleware.GenerateTokens(u.UID, "")
				if err != nil {
					ctx.JSON(
						http.StatusInternalServerError,
						gin.H{
							"error": err,
						},
					)
					return
				}
				ctx.JSON(
					http.StatusOK,
					gin.H{
						"jwt":           tokens.Jwt,
						"refresh-token": tokens.Refresh,
					},
				)
				return
			} else {
				ctx.JSON(
					http.StatusUnauthorized,
					gin.H{
						"error": "Password did not match",
					},
				)
				return
			}
		}
	}

	ctx.JSON(
		http.StatusNotFound,
		gin.H{
			"error": "user not found",
		},
	)
	return
}

func Refresh(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"err": "Authorization token not provided in header",
			},
		)
		return
	}

	auth := strings.Split(header, " ")
	if len(auth) != 2 || auth[0] != "Bearer" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"err": "Authorisation header doesn't have bearer string",
			},
		)
		return
	}

	claim, err := middleware.ParseToken(auth[1])
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"err": err,
			},
		)
		return
	}

	tokens, err := middleware.GenerateTokens(claim.UserId, "")
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err,
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"jwt":           tokens.Jwt,
			"refresh-token": tokens.Refresh,
		},
	)

}
