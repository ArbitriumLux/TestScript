package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

func main() {
	const (
		seleniumPath     = "vendor/selenium-server-standalone-3.141.59.jar"
		chromeDriverPath = "vendor/chromedriver.exe"
		port             = 8080
	)
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath),
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	if err := wd.Get("http://test.youplace.net/"); err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 1000)

	btn, err := wd.FindElement(selenium.ByTagName, "button")
	if err != nil {
		panic(err)
	}

	if err := btn.Click(); err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 1000)
	for {
		pshk, err := wd.FindElements(selenium.ByCSSSelector, "p")
		if err != nil {
			panic(err)
		}
		for i := range pshk {
			inps, err := pshk[i].FindElements(selenium.ByTagName, "input")
			if err != nil {
				println(err)
			}

			long := ""

			for i := range inps {
				s, err := inps[i].GetAttribute("type")
				if err != nil {
					panic(err)
				}

				if s == "radio" {
					attr, err := inps[i].GetAttribute("value")
					if err != nil {
						panic(err)
					}
					if len(attr) > len(long) {
						long = attr
					}
				}

				if s == "text" {
					err := inps[i].SendKeys(`test`)
					if err != nil {
						panic(err)
					}
				}
			}

			in, err := pshk[i].FindElement(selenium.ByCSSSelector, "input[value='"+long+"']")
			if err != nil {
				println(err)
			}
			if err == nil {
				if err := in.Click(); err != nil {
					println(err)
				}
			}
			sel, err := pshk[i].FindElements(selenium.ByTagName, "option")
			if err != nil {
				println(err)
			}

			longS := ""
			for i := range sel {
				attrS, err := sel[i].GetAttribute("value")
				if err != nil {
					panic(err)
				}
				if len(attrS) >= len(longS) {
					longS = attrS
				}
			}

			sl, err := pshk[i].FindElement(selenium.ByCSSSelector, "option[value='"+longS+"']")
			if err != nil {
				println(err)
			}
			if err == nil {
				if err := sl.Click(); err != nil {
					println(err)
				}
			}
		}
		time.Sleep(time.Millisecond * 3000)
		subm, err := wd.FindElement(selenium.ByTagName, "button")
		if err != nil {
			panic(err)
		}

		if err := subm.Click(); err != nil {
			panic(err)
		}
	}
}
