package datasource

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/viper"
	"github.com/starudream/go-lib/core/v2/codec/json"
	"log"
	"snowbreak_bot/utils"
)

// UpdateDataSource 更新数据源
func UpdateDataSource() func() {
	updateDataSource := func() {
		go UpdateDataSourceRunner()
	}
	return updateDataSource
}

// UpdateDataSourceRunner 更新数据源
func UpdateDataSourceRunner() {
	log.Println("开始更新数据源...")
	var characterList []utils.Character
	api := viper.GetString("api.wiki")
	pw, err := playwright.Run()
	if err != nil {
		log.Println("未检测到playwright，开始自动安装...")
		playwright.Install()
		pw, _ = playwright.Run()
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Println(err)
	}
	page, _ := browser.NewPage()
	defer func() {
		page.Close()
	}()
	page.Goto(api+"/snow", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	locator, _ := page.Locator(".more")
	more, _ := locator.Nth(1)
	err = more.Click()
	if err != nil {
		log.Println(err)
		return
	}
	page.WaitForTimeout(3000)
	html, _ := page.Content()
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
	if err != nil {
		return
	}
	doc.Find(".title").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() == "角色图鉴" {
			selection.Parent().Next().Find(".item-wrapper").Eq(0).Find("a").Each(func(j int, selection *goquery.Selection) {
				n := selection.Text()
				var char utils.Character
				char.Name = n
				for _, attr := range selection.Nodes[0].FirstChild.NextSibling.Attr {
					if attr.Key == "src" {
						char.ThumbURL = "https:" + attr.Val
					}
				}
				characterList = append(characterList, char)
			})
		}
	})

	utils.RedisSet("characterList", json.MustMarshalString(characterList), 0)

	// 武器
	/*var weaponList []utils.Weapon
	pw, err := playwright.Run()
	if err != nil {
		log.Println("未检测到playwright，开始自动安装...")
		playwright.Install()
		pw, _ = playwright.Run()
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Println(err)
	}
	page, _ := browser.NewPage()
	defer func() {
		page.Close()
	}()
	page.Goto(api+"/snow", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	page.Click(".primary-btn")
	page.WaitForTimeout(1000)
	page.Click(".more")
	page.WaitForTimeout(3000)
	html, _ := page.Content()
	doc, err = goquery.NewDocumentFromReader(bytes.NewBufferString(html))
	if err != nil {
		return
	}
	doc.Find(".title").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() == "武器图鉴" {
			selection.Parent().Next().Find(".item-wrapper").Eq(0).Find("a").Each(func(j int, selection *goquery.Selection) {
				var weapon utils.Weapon
				weapon.Name = selection.Text()
				href, _ := selection.Attr("href")
				weapon.Url = api + href
				for c, attr := range selection.Nodes[0].FirstChild.NextSibling.Attr {
					if attr.Key == "data-src" {
						weapon.ThumbURL = "https:" + selection.Nodes[0].FirstChild.NextSibling.Attr[c].Val
						weaponList = append(weaponList, weapon)
						break
					}
				}
			})
		}
	})

	utils.RedisSet("weaponList", json.MustMarshalString(weaponList), 0)*/

	log.Println("数据源更新完毕")
}
