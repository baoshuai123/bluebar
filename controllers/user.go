package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignupHandler 处理注册请求函数
func SignupHandler(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamSignUP)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数错误,直接返回响应
		zap.L().Error("SignUP with invalid param", zap.Error(err))
		//判断 err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)), //翻译错误
		//})
		return
	}
	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error(" logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success",
	//})
}

// LoginHandler 用户登录函数
func LoginHandler(c *gin.Context) {
	//1.获取请求参数参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数错误,直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断 err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)), //翻译错误
		//})
		return
	}
	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "账号密码错误",
		//})
	}
	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})

}
