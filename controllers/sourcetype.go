package controllers
import (
	"github.com/astaxie/beego/orm"
	"common.dh.cn/utils"
	"common.dh.cn/models"
	"common.dh.cn/controllers"
	"fmt"
	"strconv"
)
type SourcetypeController struct {
	controllers.BaseController
}
func (c *SourcetypeController) init(i int) {
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
	c.Data["Menu"]=Menu
	Authname:=c.Ctx.GetCookie("Authname")
	c.Data["Authname"]=Authname
}


func (c *SourcetypeController) List() {
	c.init(2)
	var mpurl ="/sourcetype/list?"
	c.TplName = "sourcetype/index.html"
	var total,total_page int64
	var list []*models.DiDatasourceType
	page,_ := c.GetInt64("page",1)
	page_size,_ := c.GetInt64("page_size",10)
	search := c.GetString("search")
	status:= c.GetString("status")
	filters := map[string]interface{}{}
	if len(search)>0{
		cond := orm.NewCondition()

		if len(search)>0{
			c.Data["search"] = search
			mpurl=mpurl+"&search="+search
			condor := cond.Or("name__icontains", search)

			if len(status)>0{
				c.Data["status"] = status
				int,_:=strconv.Atoi(status)
				mpurl=mpurl+"&status="+status
				condor=cond.AndCond(condor).And("status",int)
			}else {

				c.Data["status"] = "nil"
			}

			number,_:=new(models.DiDatasourceType).Query().Offset((page-1)*page_size).Limit(page_size).SetCond(condor).OrderBy("-create_time").All(&list)
			total,_=new(models.DiDatasourceType).Query().SetCond(condor).Count()
			if total%page_size!=0{
				total_page=total/page_size+1
			}else {
				total_page=total/page_size
			}


			fmt.Println(number)



		}

	}else {
		if len(status)>0{
			c.Data["status"] = status
			int,_:=strconv.Atoi(status)
			filters["status"]=int
			mpurl=mpurl+"&status="+status
		}else {

			c.Data["status"] = "nil"
		}



		total,total_page,list = new(models.DiDatasourceType).OrderPager(page, page_size, filters,"-create_time")
	}
	data := []utils.P{}
	if len(list) > 0 {
		for _, info := range list {
			dhdatasource := utils.P{}
			dhdataType := utils.P{}
			dhdatasource["datasource_type_id"]=info.ObjectId
			dhdatasourceCount:=new(models.DiDatasourcePub).Count(dhdatasource)
			dhdataType["ObjectId"] = info.ObjectId
			dhdataType["Name"] = info.Name
			dhdataType["SourceCount"] = dhdatasourceCount
			dhdataType["CreateTime"] = info.CreateTime.Format("2006-01-02 15:04:05")
			dhdataType["Status"] = info.Status
			data = append(data, dhdataType)
		}
	}
	c.Data["List"] = data
	c.Data["Pagination"] = PagerHtml(int(total), int(total_page), int(page_size), int(page),mpurl)
}
func (c *SourcetypeController) Update() {
	c.Require("id")
	id := c.GetString("id")
	name := c.GetString("name")
	DiDatasourceType := new(models.DiDatasourceType).Find(id)
	if DiDatasourceType == nil {
		c.EchoJsonErr("分类组不存在")
		c.StopRun()
	}
	if name!=""{
		DiDatasourceType.Name=name
	}
	if c.GetString("status") == "1" {
		int,_:=strconv.Atoi(c.GetString("status"))
		DiDatasourceType.Status = int
	}
	if c.GetString("status") == "0"{
		int,_:=strconv.Atoi(c.GetString("status"))
		DiDatasourceType.Status = int
	}
	fmt.Println("---------------a-------------------",DiDatasourceType)
	result := DiDatasourceType.Save()
	if !result {
		c.EchoJsonErr("修改失败")
	} else {
		c.EchoJsonOk()
	}
}
func (c *SourcetypeController) Remove() {
	c.Require("id")
	id := c.GetString("id")
	DiDatasourceType := new(models.DiDatasourceType).Find(id)
	if DiDatasourceType == nil {
		c.EchoJsonErr("数据源组不存在")
		c.StopRun()
	}
	result := DiDatasourceType.Delete(id)
	if !result {
		c.EchoJsonErr("删除失败")
	} else {
		c.EchoJsonOk()
	}
}
func (c *SourcetypeController) Edit() {
	c.Require("id")
	id := c.GetString("id")
	DiDatasourceType := new(models.DiDatasourceType).Find(id)
	if DiDatasourceType == nil {
		c.EchoJsonErr("用户不存在")
		c.StopRun()
	}

	c.Data["object"] = &DiDatasourceType
	c.TplName = "sourcetype/edit.html"
}
func (c *SourcetypeController) Add() {
	DiDatasourceType := new(models.DiDatasourceType)
	DiDatasourceType.Name = c.GetString("name")
	DiDatasourceType.Status = 0
	result :=DiDatasourceType.Save()
	if !result {
		c.EchoJsonErr("创建失败")
	} else {
		c.EchoJsonOk()
	}
}
func (c *SourcetypeController) Listremove() {
	c.Require("datas")
	datas := c.GetString("datas")
	plist:=*utils.JsonDecodeArrays([]byte(datas))
	argerr :=make([]string,1)
	for _,v :=range plist{
		DiDatasourceType := new(models.DiDatasourceType).Find(v["object_id"].(string))
		if DiDatasourceType == nil {
			argerr=append(argerr,v["object_id"].(string))
		}else {
			result := DiDatasourceType.Delete(v["object_id"].(string))
			if !result {
				argerr=append(argerr,v["object_id"].(string))
			}
		}
	}
	if    len(argerr[0]) > 0 {
		c.EchoJsonErr("部分删除失败")
	}
	c.EchoJsonOk()


}

func (c *SourcetypeController) Create() {
	c.TplName = "sourcetype/create.html"
}