package handler

import (
	"net/http"
	"rest-api/rest-api/user"
	"rest-api/rest-api/utils"

	"github.com/gin-gonic/gin"
)

var users []user.User

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
	user := user.User{}
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

	userUpdateReq := user.UserUpdateReq{}
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

	userUpdateReq := user.UserUpdateReq{}
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
