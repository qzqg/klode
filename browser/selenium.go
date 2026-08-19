package browser

import (
	"fmt"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func SetupBrowser() (selenium.WebDriver, error) {
	// إعداد خيارات Chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--disable-blink-features=AutomationControlled",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)

	// الاتصال بـ ChromeDriver (بدون /wd/hub)
	wd, err := selenium.NewRemote(caps, "http://localhost:4444")
	if err != nil {
		return nil, fmt.Errorf("فشل في الاتصال بـ ChromeDriver: %v", err)
	}

	return wd, nil
}

func LoginToChess(wd selenium.WebDriver, email, password string) error {
	if err := wd.Get("https://www.chess.com/login"); err != nil {
		return fmt.Errorf("فشل في فتح صفحة تسجيل الدخول: %v", err)
	}
	time.Sleep(3 * time.Second)

	emailInput, err := wd.FindElement(selenium.ByName, "username")
	if err != nil {
		emailInput, err = wd.FindElement(selenium.ByCSSSelector, "input[name='username']")
		if err != nil {
			return fmt.Errorf("لم يتم العثور على حقل البريد الإلكتروني: %v", err)
		}
	}
	emailInput.Clear()
	emailInput.SendKeys(email)

	passInput, err := wd.FindElement(selenium.ByName, "password")
	if err != nil {
		passInput, err = wd.FindElement(selenium.ByCSSSelector, "input[name='password']")
		if err != nil {
			return fmt.Errorf("لم يتم العثور على حقل كلمة المرور: %v", err)
		}
	}
	passInput.Clear()
	passInput.SendKeys(password)

	loginBtn, err := wd.FindElement(selenium.ByCSSSelector, "button[type='submit']")
	if err != nil {
		return fmt.Errorf("لم يتم العثور على زر تسجيل الدخول: %v", err)
	}
	loginBtn.Click()

	time.Sleep(3 * time.Second)
	return nil
}

func GetFENFromBoard(wd selenium.WebDriver) (string, error) {
	boardElement, err := wd.FindElement(selenium.ByID, "board")
	if err != nil {
		return "", fmt.Errorf("لم يتم العثور على اللوحة: %v", err)
	}

	fen, err := boardElement.GetAttribute("data-fen")
	if err != nil {
		return "", fmt.Errorf("فشل في قراءة FEN: %v", err)
	}

	if fen == "" {
		return "", fmt.Errorf("لم يتم العثور على FEN (اللوحة فارغة؟)")
	}

	return fen, nil
}

func ExecuteMove(wd selenium.WebDriver, move string) error {
	from := move[:2]
	to := move[2:4]

	fmt.Printf("🔄 تنفيذ النقلة من %s إلى %s\n", from, to)

	// هنا يمكنك إضافة منطق الضغط على المربعات
	return nil
}