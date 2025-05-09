package helper

import (
	"encoding/json"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
)

type user struct {
	ID int64 `json:"id"`
}

func GetUserByToken(ctx *gin.Context) (int64, error) {
	claims, exist := ctx.Get("claims")
	if !exist {
		return 0, nil
	}
	var user user
	err := json.Unmarshal(claims.([]byte), &user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// function ถ้าจะใช่ให้ไปสร้าง cloudinay account ก่อน
func NewCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, err
	}
	return cld, nil
}
