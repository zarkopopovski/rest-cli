package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type AuthView struct {
	InputKey        bool
	InputValue      bool
	OptSelector     bool
	IsKeyUsername   bool
	IsValuePassword bool
	OptSelectorData []string
	SelectedIndex   int
	AuthDataMap     map[string]string
	AuthKeyLoc      string
	AkSelect        *widget.Select
}

func (self *AuthView) InitForm(inputKey bool, inputValue bool, optSelector bool, optSelectorData []string, selectedIndex int, isKeyUsername bool, isValuePassword bool) {
	self.InputKey = inputKey
	self.InputValue = inputValue
	self.OptSelector = optSelector
	self.IsKeyUsername = isKeyUsername
	self.IsValuePassword = isValuePassword
	self.OptSelectorData = optSelectorData
	self.SelectedIndex = selectedIndex
}

func (self *AuthView) BuildAuthForm(myWindow fyne.Window) *widget.Form {
	formItems := []*widget.FormItem{}
	var akInputKey *widget.Entry
	var akInputValue *widget.Entry
	//var akSelect *widget.Select

	if self.InputKey == true {
		akInputKey = widget.NewEntry()
		if self.IsKeyUsername == true {
			akInputKey.SetPlaceHolder("Enter username...")
		} else {
			akInputKey.SetPlaceHolder("Enter key...")
		}
		formItems = append(formItems, &widget.FormItem{Text: "Key", Widget: akInputKey})
	}
	if self.InputValue == true {
		if self.IsValuePassword == true {
			akInputValue = widget.NewPasswordEntry()
			akInputValue.SetPlaceHolder("Enter password...")
		} else {
			akInputValue = widget.NewEntry()
			akInputValue.SetPlaceHolder("Enter value...")
		}
		formItems = append(formItems, &widget.FormItem{Text: "Value", Widget: akInputValue})
	}
	if self.OptSelector == true {
		self.AuthKeyLoc = self.OptSelectorData[self.SelectedIndex]
		self.AkSelect = widget.NewSelect(self.OptSelectorData, func(value string) {
			self.AuthKeyLoc = value
		})
		self.AkSelect.SetSelected(self.AuthKeyLoc)
		formItems = append(formItems, &widget.FormItem{Text: "Add top", Widget: self.AkSelect})
	}

	akForm := &widget.Form{
		Items: formItems, OnSubmit: func() {
			self.AuthDataMap = make(map[string]string)
			if self.OptSelector == true {
				self.AuthDataMap["authKeyLoc"] = self.AuthKeyLoc
			}
			if self.InputKey == true {
				self.AuthDataMap["authKey"] = akInputKey.Text
			}
			if self.InputValue == true {
				self.AuthDataMap["authValue"] = akInputValue.Text
			}

			dialog.ShowInformation("Auth Info", "Params are set", myWindow)
		}, SubmitText: "Add Key",
	}
	akForm.Hide()

	return akForm
}
