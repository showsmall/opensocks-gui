package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/net-byte/opensocks/client"
	"github.com/net-byte/opensocks/config"
)

func main() {
	app := app.New()
	win := app.NewWindow("openscoks-gui")
	win.Resize(fyne.NewSize(320, 150))
	config := loadConfig()
	localAddr := widget.NewEntry()
	localAddr.Text = config.LocalAddr
	serverAddr := widget.NewEntry()
	serverAddr.Text = config.ServerAddr
	username := widget.NewEntry()
	username.Text = config.Username
	password := widget.NewPasswordEntry()
	password.Text = config.Password

	appName := widget.NewLabelWithStyle("OpenSocks v1.0.0", fyne.TextAlignCenter, fyne.TextStyle{})
	msg := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("local addr:", localAddr),
			widget.NewFormItem("server addr:", serverAddr),
			widget.NewFormItem("username:", username),
			widget.NewFormItem("password:", password),
		},
	}
	connectBtn := widget.NewButtonWithIcon("connect", theme.MailSendIcon(), func() {
		var err error
		config.LocalAddr = localAddr.Text
		config.ServerAddr = serverAddr.Text
		config.Username = username.Text
		config.Password = password.Text
		config.Wss = true
		if config.LocalAddr == "" || config.ServerAddr == "" {
			msg.Text = "addr can't be empty!"
			return
		}
		_, err = net.ResolveTCPAddr("tcp", config.LocalAddr)
		if nil != err {
			msg.Text = "error local addr format!"
			return
		}
		_, err = net.ResolveTCPAddr("tcp", config.ServerAddr)
		if nil != err {
			msg.Text = "error server addr format!"
			return
		}
		go client.Start(config)
		msg.Text = "successfully connected!"
		saveConfig(config)
	})
	exitBtn := widget.NewButtonWithIcon("exit", theme.CancelIcon(), func() {
		win.Close()
	})

	box := widget.NewVBox(appName, form, connectBtn, exitBtn, msg)

	win.SetContent(box)
	win.ShowAndRun()
}

func loadConfig() config.Config {
	var result config.Config
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		log.Println(err)
		//set default config
		result = config.Config{}
		result.LocalAddr = "127.0.0.1:1081"
		result.ServerAddr = "example.com:443"
		result.Username = "admin"
		result.Password = "pass@123456"
		result.Wss = true
		return result
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func saveConfig(config config.Config) {
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile("./config.json", file, 0644)
}
