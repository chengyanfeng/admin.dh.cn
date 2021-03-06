package controllers

import (
	"common.dh.cn/models"
	"common.dh.cn/controllers"
	"common.dh.cn/utils"
	"fmt"


)

type LoginController struct {

	controllers.BaseController
}
func (c *LoginController) init(i int) {
c.Layout = "common/layout.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "common/header.html"
	c.LayoutSections["HtmlFooter"] = "common/footer.html"
	for   k,v :=range Menu{
		if k!=i{
			v["On"]=0
		}else {
			Menu[i]["On"]=1
			if Menu[i]["Sub"]!=nil{
				a:= Menu[i]["Sub"].(interface{})
				b:= a.([]utils.P)
				for _,v:=range b {
					v["On"]=1
				}
			}
		}
	}
	Authname:=c.Ctx.GetCookie("Authname")
	c.Data["Authname"]=Authname
	c.Data["Menu"]=Menu
}

func (c *LoginController)Get(){
	utils.SDel("gooid")
	c.SetSecureCookie("","gooid","")
	c.TplName="index/login.html"
}

func (c *LoginController)Login(){
	fmt.Println("---------------------------dafdsffdsa---------")
	c.init(0)
	username:=c.GetString("name")
	password:=c.GetString("password")
	userfilter:=map[string]interface{}{}
	userfilter["name"]=username
	user:=new(models.DhUser).Find(userfilter)
	fmt.Print(user)
	if user==nil {
		c.EchoJsonErr("用户不存在")
		c.StopRun()
	}
	if password  != user.Password{

		c.EchoJsonErr("密码不正确")
		c.StopRun()
	}
	fmt.Println("111111111111111women1111111111111111")


	c.SetSession("gooid",username)
	utils.S("gooid",username,"7200")
	c.Ctx.SetCookie("Authname",username)
	/*

	c.Ctx.SetCookie("gooid",username+utils.ToString(time.Second))
	*/
	c.EchoJsonOk("index")
}
func (c *LoginController)Quit(){
	c.DelSession("gooid")
	c.Redirect("/",200)

}