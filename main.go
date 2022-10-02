package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/canvas"

	"fyne.io/fyne/v2/layout"
	
	//"encoding/json"
	//"fmt"
	//"image/color"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"net/http/cookiejar"
	"errors"
	"bytes"
	"os"
	"strings"
	"strconv"
	"log"
  "mime/multipart"
	//"path/filepath"  
	
	"github.com/zarkopopovski/rest-cli/services"
	"github.com/zarkopopovski/rest-cli/models"
)

type RestClient struct {
	pInputURL 		string
	httpMethod		string
	radioOpt1 		string
	radioOpt2 		string
	paramTextV		string
	authMethod		string
	authKeyLoc		string
	selectedFile	string
	pDataMap			map[string]string
	dataMap		 	map[string]string
	cDataMap			map[string]string
	bDataMap    	map[string]string
	authDataMap   map[string]string
	selColIDX	   int
	selColID 	   int
	selReqIDX	   int
	selReqID 	   int
	collections   []*models.Collection
	colRequests 	[]*models.UrlRequest
	client 		 	*http.Client	
	req 					*http.Request
	res 					*http.Response
	myWindow			fyne.Window
	DBService	 	*services.DBService
}

func (self *RestClient) BuildAuthForm(inputKey bool, inputValue bool, optSelector bool, optSelectorData []string, selectedIndex int, isKeyUsername bool, isValuePassword bool) (*widget.Form) {
	formItems := []*widget.FormItem{}
	var akInputKey *widget.Entry
	var akInputValue *widget.Entry
	var akSelect *widget.Select
	
	if inputKey == true {
		akInputKey = widget.NewEntry()
		if isKeyUsername == true {
			akInputKey.SetPlaceHolder("Enter username...")
		} else {
			akInputKey.SetPlaceHolder("Enter key...")
		}
		formItems = append(formItems, &widget.FormItem{Text:"Key", Widget:akInputKey})
	}
	if inputValue == true {
		if isValuePassword == true {
			akInputValue = widget.NewPasswordEntry()
			akInputValue.SetPlaceHolder("Enter password...")
		} else {
			akInputValue = widget.NewEntry()
			akInputValue.SetPlaceHolder("Enter value...")
		}
		formItems = append(formItems, &widget.FormItem{Text:"Value", Widget:akInputValue})
	}
	if optSelector == true {
		self.authKeyLoc = optSelectorData[selectedIndex]
		akSelect = widget.NewSelect(optSelectorData, func(value string) {
			self.authKeyLoc = value
		})
		akSelect.SetSelected(self.authKeyLoc)
		formItems = append(formItems, &widget.FormItem{Text:"Add top", Widget:akSelect})
	}
	
	akForm := &widget.Form{
		Items: formItems, OnSubmit: func() {
			self.authDataMap = make(map[string]string)
			if optSelector == true {
				self.authDataMap["authKeyLoc"] = self.authKeyLoc
			}
			if inputKey == true {
				self.authDataMap["authKey"] = akInputKey.Text
				akInputKey.SetText("")
			}
			if inputValue == true {
				self.authDataMap["authValue"] = akInputValue.Text
				akInputValue.SetText("")
			}
			
			dialog.ShowInformation("Auth Info", "Params are set", self.myWindow)
		}, SubmitText: "Add Key",
	}
	akForm.Hide()
	
	return akForm
}

func (self *RestClient) ResetData() {
	self.pInputURL = ""
	self.httpMethod = ""
	self.radioOpt1 = ""
	self.radioOpt2 = ""
	self.paramTextV = ""
	self.authMethod = ""
	self.authKeyLoc = ""
	self.selectedFile = ""
	self.pDataMap = nil
	self.dataMap = nil
	self.cDataMap = nil
	self.bDataMap = nil
	self.authDataMap = nil
}

func (self *RestClient) BuildUI() {
	myApp := app.New()
	self.myWindow = myApp.NewWindow("Rest Client")
	//self.myWindow.SetFixedSize(true)
	
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter URL...")
	
	largeText := widget.NewMultiLineEntry()	
	largeText.SetPlaceHolder("Response...")
	largeText.Resize(fyne.NewSize(390,564))
	
	combo := widget.NewSelect([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"}, func(value string) {
		self.httpMethod = value
	})
	
	//PARAMS TAB START
	prmInputKey := widget.NewEntry()
	prmInputKey.SetPlaceHolder("Enter param key...")
	
	prmInputValue := widget.NewEntry()
	prmInputValue.SetPlaceHolder("Enter param value...")
	
	self.pDataMap = make(map[string]string, 1)
	pData := binding.BindStringList(
		&[]string{},
	)
	pList := widget.NewListWithData(pData,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
		
	pList.Resize(fyne.NewSize(400, 200))
	
	pListAreaContent := container.NewWithoutLayout(pList)
	
	pForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Key", Widget: prmInputKey}, 
			{Text: "Value", Widget: prmInputValue},
		}, OnSubmit: func() {
			pData.Append(prmInputKey.Text + " | " + prmInputValue.Text)
			self.pDataMap[prmInputKey.Text] = prmInputValue.Text
			
			if len(self.pDataMap) > 0 {
				if self.pInputURL == "" {
					self.pInputURL = input.Text
				}
				
				tempURL := self.pInputURL
				
				urlParams := url.Values{}
				
				for key, value := range self.pDataMap {
					urlParams.Add(key, value)
				}
				
				tempURL = tempURL+"?"+urlParams.Encode()
				
				input.SetText(tempURL)				
			}
			
			prmInputKey.SetText("")
			prmInputValue.SetText("")
		}, SubmitText: "Add Key",
	}
	//PARAMS TAB END
	
	//HEADER TAB START
	hdrInputKey := widget.NewEntry()
	hdrInputKey.SetPlaceHolder("Enter header key...")
	
	hdrInputValue := widget.NewEntry()
	hdrInputValue.SetPlaceHolder("Enter header value...")
	
	self.dataMap = make(map[string]string, 1)
	data := binding.BindStringList(
		&[]string{},
	)
	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
		
	list.Resize(fyne.NewSize(400, 200))
	
	listAreaContent := container.NewWithoutLayout(list)
	
	hForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Key", Widget: hdrInputKey}, 
			{Text: "Value", Widget: hdrInputValue},
		}, OnSubmit: func() {
			data.Append(hdrInputKey.Text + " | " + hdrInputValue.Text)
			self.dataMap[hdrInputKey.Text] = hdrInputValue.Text
			
			hdrInputKey.SetText("")
			hdrInputValue.SetText("")
		}, SubmitText: "Add Key",
	}
	//HEADER TAB END

	//COOKIE TAB START
	cooInputKey := widget.NewEntry()
	cooInputKey.SetPlaceHolder("Enter cookie key...")
	
	cooInputValue := widget.NewEntry()
	cooInputValue.SetPlaceHolder("Enter cookie value...")
	
	self.cDataMap = make(map[string]string, 1)
	cData := binding.BindStringList(
		&[]string{},
	)
	cList := widget.NewListWithData(cData,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
		
	cList.Resize(fyne.NewSize(400, 200))
	
	cListAreaContent := container.NewWithoutLayout(cList)
	
	cForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Key", Widget: cooInputKey}, 
			{Text: "Value", Widget: cooInputValue},
		}, OnSubmit: func() {
			cData.Append(cooInputKey.Text + " | " + cooInputValue.Text)
			self.cDataMap[cooInputKey.Text] = cooInputValue.Text
			
			cooInputKey.SetText("")
			cooInputValue.SetText("")
		}, SubmitText: "Add Key",
	}
	//COOKIE TAB END

	//BODY TAB START
	var bForm = &widget.Form{}
		
	booInputKey := widget.NewEntry()
	booInputKey.SetPlaceHolder("Enter body key...")
	
	booInputValue := widget.NewEntry()
	booInputValue.SetPlaceHolder("Enter body value...")
	
	self.bDataMap = make(map[string]string, 1)
	bData := binding.BindStringList(
		&[]string{},
	)
	bList := widget.NewListWithData(bData,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
		
	bList.Resize(fyne.NewSize(400, 200))
	
	bListAreaContent := container.NewWithoutLayout(bList)
	bListAreaContent.Show()
	
	paramText := widget.NewMultiLineEntry()	
	paramText.SetPlaceHolder("Body Params...")
	paramText.Resize(fyne.NewSize(400,300))
	paramText.Hide()
	
	comboOpt2 := widget.NewSelect([]string{"Text", "JavaScript", "JSON", "HTML", "XML"}, func(value string) {
		self.radioOpt2 = value
		if self.radioOpt2 == "json" {
			
		}
	})
	comboOpt2.Disable()
	
	fileInput := widget.NewEntry()
	fileInput.Disable()
	fileInputLabel := widget.NewLabel("File     ")
	fileContentBorder := container.New(layout.NewBorderLayout(nil, nil, fileInputLabel, nil), fileInputLabel, fileInput)	
	fileOpenDialog := widget.NewButton("File Open", func() {
		dialog.ShowFileOpen(func(read fyne.URIReadCloser, err error) {
      if err != nil {
        dialog.ShowError(err, self.myWindow)
        return
      }
      if read == nil {
        return
      }
      selectedFile := read.URI().String()
      filePath := strings.TrimPrefix(selectedFile, "file://")
      fileInput.SetText(filePath)
      self.selectedFile = filePath
    }, self.myWindow)		
	})
	fileVBox := container.NewVBox(fileContentBorder, fileOpenDialog)
	fileVBox.Hide()
		
	bodyTypeLabel := widget.NewLabel("Type  ")
	comboOpt1 := widget.NewSelect([]string{"none", "form-data", "x-www-form-urlencoded", "raw", "binary"}, func(value string) {
		self.radioOpt1 = value
		if self.radioOpt1 == "raw" {
			comboOpt2.Enable()
			bListAreaContent.Hide()
			paramText.Show()
			bForm.Disable()
			booInputKey.Disable()
			booInputValue.Disable()
			fileVBox.Hide()
			self.selectedFile = ""
		} else if self.radioOpt1 == "binary" {
			bForm.Hide()
			bListAreaContent.Hide()
			paramText.Hide()
			fileVBox.Show()
			self.selectedFile = ""
		} else if self.radioOpt1 == "form-data" {
			comboOpt2.Disable()
			bListAreaContent.Show()
			paramText.Hide()
			bForm.Show()
			bForm.Enable()
			booInputKey.Enable()
			booInputValue.Enable()
			fileVBox.Show()
			self.selectedFile = ""
		} else {
			comboOpt2.Disable()
			bListAreaContent.Show()
			paramText.Hide()
			bForm.Enable()
			booInputKey.Enable()
			booInputValue.Enable()
			fileVBox.Hide()
			self.selectedFile = ""
		}		
	})	
	
	bForm = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Key", Widget: booInputKey}, 
			{Text: "Value", Widget: booInputValue},
			{Text: "Sub", Widget: comboOpt2},
		}, OnSubmit: func() {
			bData.Append(booInputKey.Text + " | " + booInputValue.Text)
			self.bDataMap[booInputKey.Text] = booInputValue.Text
			
			booInputKey.SetText("")
			booInputValue.SetText("")
		}, SubmitText: "Add Key",
	}
	//BODY TAB END
	
	akForm := &widget.Form{}
	brtForm := &widget.Form{}
	basForm := &widget.Form{}
	
	self.authMethod = "No Auth"
	authLabel := widget.NewLabel("Auth Type")
	authSelectOpt := []string{"No Auth", "API Key", "Bearer Token", "Basic Auth"}
	authSelect := widget.NewSelect(authSelectOpt, func(value string) {
		self.authMethod = value
		if value == authSelectOpt[0] {
			akForm.Hide()
			brtForm.Hide()
			basForm.Hide()
		} else if value == authSelectOpt[1] {
			akForm.Show()
			brtForm.Hide()
			basForm.Hide()
		} else if value == authSelectOpt[2] {
			akForm.Hide()
			brtForm.Show()
			basForm.Hide()
		} else if value == authSelectOpt[3] {
			akForm.Hide()
			brtForm.Hide()
			basForm.Show()
		}
	})
	authSelect.SetSelected(self.authMethod)	
	
	authBorder := container.New(layout.NewBorderLayout(nil, nil, authLabel, nil), authLabel, authSelect)	
	
	bodyBorder := container.New(layout.NewBorderLayout(nil, nil, bodyTypeLabel, nil), bodyTypeLabel, comboOpt1)	
	
	akForm = self.BuildAuthForm(true, true, true, []string{"Header", "Query Params"}, 0, false, false)
	brtForm = self.BuildAuthForm(false, true, false, []string{}, 0, false, false)
	basForm = self.BuildAuthForm(true, true, false, []string{}, 0, true, true)

	formSeparator := widget.NewSeparator()
		
	tabs := container.NewAppTabs(
		container.NewTabItem("Params", container.NewVBox(pForm, pListAreaContent)),
		container.NewTabItem("Header", container.NewVBox(hForm, listAreaContent)),
		container.NewTabItem("Cookie", container.NewVBox(cForm, cListAreaContent)),
		container.NewTabItem("Body", container.NewVBox(bodyBorder, bForm, fileVBox, bListAreaContent, paramText)),
		container.NewTabItem("Auth", container.NewVBox(authBorder, formSeparator, akForm, brtForm, basForm)),
	)
	
	form := &widget.Form{
		Items: []*widget.FormItem{
		{Text: "URL:", Widget: input}, 
		{Text: "METHOD:", Widget: combo},
		//{Text: "OPTION:", Widget: wrapCheck}
		}, OnSubmit: func() { 
			if (input.Text == "") {
				dialog.ShowError(errors.New("Empty URL"), self.myWindow)
				return
			}
			
			if (self.httpMethod == "") {
				dialog.ShowError(errors.New("HTTP method not selected"), self.myWindow)
				return
			}
			
			urlString := input.Text
			self.paramTextV = paramText.Text
			
			self.ExecuteRequest(urlString, func(stringRes string) {
				largeText.SetText(stringRes)
			})			
		},
	}
	
	contentV1 := container.NewVBox(form, tabs)
	
	grid := container.New(layout.NewGridLayout(2), contentV1, largeText)
	
	collections, err := self.DBService.ListAllCollection()
	if err != nil {
		dialog.ShowError(errors.New("Collection loading problem or there are not any saved"), self.myWindow)
		return
	}
	self.collections = collections
	
	reqList := &widget.List{} 
	 
	colList := widget.NewList(
		func() int { return len(self.collections) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
				co.(*widget.Label).SetText(self.collections[lii].Name)
		},
	)
	colList.OnSelected = func(id widget.ListItemID) {
		selCollectionObj := self.collections[id]
		self.selColID  = int(selCollectionObj.Id)
		self.selColIDX = id
		colRequests, err := self.DBService.ListAllRequestsForCollection(selCollectionObj.Id)
		if err != nil {
			dialog.ShowError(errors.New("Collection requests loading problem or there are not any saved"), self.myWindow)
			return
		}
		if len(colRequests) > 0 {
			self.colRequests = colRequests
		} else {
			self.colRequests = nil
		}	
		reqList.Refresh()
	}
	colList.Resize(fyne.NewSize(100, 50))
	
	reqList = widget.NewList(
		func() int { return len(self.colRequests) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			if len(self.colRequests) > 0 {
				co.(*widget.Label).SetText(self.colRequests[lii].Method + " " + self.colRequests[lii].Url)
			}
		},
	)
	reqList.OnSelected = func(id widget.ListItemID) {
		selRequestObj := self.colRequests[id]
		self.selReqIDX = id
		self.selReqID = int(selRequestObj.Id)
		
		self.pInputURL = selRequestObj.Url
		input.SetText(selRequestObj.Url)
		
		self.httpMethod = selRequestObj.Method
		combo.SetSelected(selRequestObj.Method)
	}
	reqList.Resize(fyne.NewSize(100, 400))

	var leftContent fyne.CanvasObject
	
	collectionInput := widget.NewEntry()
	collectionInput.SetPlaceHolder("Collection name...")
	btnConfirmCollection := widget.NewButton("Confirm", func() {
		if collectionInput.Text != "" {
			err := self.DBService.CreateNewCollection(collectionInput.Text)
			if err != nil {
				dialog.ShowError(errors.New("Collection saving problem"), self.myWindow)
				return
			}
		}
		collections, err := self.DBService.ListAllCollection()
		if err != nil {
			dialog.ShowError(errors.New("Collection loading problem"), self.myWindow)
			return
		}
		self.collections = collections
		
		colList.Refresh()
		
		collectionInput.SetText("")
	})
		
	leftContent = container.NewVBox(collectionInput, btnConfirmCollection, colList)
	
	leftContentBorder := container.New(layout.NewBorderLayout(leftContent, nil, nil, nil), leftContent, reqList)	
	
	menuItem1 := fyne.NewMenuItem("New Request", 
		func() { 
			input.SetText("") 
			largeText.SetText("")
			
			self.ResetData()
		},
	)
	menuItem2 := fyne.NewMenuItem("Save Request", 
		func() { 
			err := self.DBService.CreateNewCollectionRequest(int64(self.selColID), map[string]string{
				"name": self.httpMethod + " " + input.Text, 
				"url": input.Text, 
				"method": self.httpMethod, 
				"params_data": "",
				"header_data": "",
				"cookie_data": "",
				"body_data": "",
			})
			if err != nil {
				dialog.ShowError(errors.New("Error saving request"), self.myWindow)
				return
			}
			dialog.ShowInformation("Successful", "Selected request was successfully deleted", self.myWindow)
			
			colRequests, err := self.DBService.ListAllRequestsForCollection(int64(self.selColID))
			if err != nil {
				dialog.ShowError(errors.New("Collection requests loading problem or there are not any saved"), self.myWindow)
				return
			}
			if len(colRequests) > 0 {
				self.colRequests = colRequests
			}	
			
			reqList.UnselectAll()
			reqList.Refresh()
			
			input.SetText("") 
			largeText.SetText("")
		},
	)
	menuItem3 := fyne.NewMenuItem("Delete Request", 
		func() { 
			dialog.ShowConfirm("Confirm", "Are you sure you want to delete this request?", func(confirm bool) {
				if confirm == true {
					if self.selReqID >= 0 {
						err := self.DBService.DeleteCollectionRequest(int64(self.selReqID))
						if err != nil {
							dialog.ShowError(errors.New("Error deleting request"), self.myWindow)
							return
						}
						dialog.ShowInformation("Successful", "Selected request was successfully deleted", self.myWindow)
						
						colRequests, err := self.DBService.ListAllRequestsForCollection(int64(self.selColID))
						if err != nil {
							dialog.ShowError(errors.New("Collection requests loading problem or there are not any saved"), self.myWindow)
							return
						}
						if len(colRequests) > 0 {
							self.colRequests = colRequests
						} else {
							self.colRequests = nil
						}	
						
						reqList.UnselectAll()
						reqList.Refresh()
						
						self.selColIDX = 0
						self.selReqIDX = 0
					}
				}
			}, self.myWindow)
		},
	)
	menuItem4 := fyne.NewMenuItem("Delete Collection", 
		func() { 
			dialog.ShowConfirm("Confirm", "Are you sure you want to delete this collection?", func(confirm bool) {
				if confirm == true {
					err := self.DBService.DeleteCollection(int64(self.selColID))
					if err != nil {
						dialog.ShowError(errors.New("Error deleting collection"), self.myWindow)
						return
					}
					dialog.ShowInformation("Successful", "Selected request was successfully deleted", self.myWindow)
						
					collections, err := self.DBService.ListAllCollection()
					if err != nil {
						dialog.ShowError(errors.New("Collection loading problem"), self.myWindow)
						return
					}
					self.collections = collections
					
					colList.UnselectAll()
					reqList.UnselectAll()
					
					self.colRequests = nil
					
					reqList.Refresh()
					
					self.selColIDX = 0
				}
			}, self.myWindow)			
		},
	)
  newMenu := fyne.NewMenu("File", menuItem1, menuItem2, menuItem3, menuItem4)
  menu := fyne.NewMainMenu(newMenu)
	self.myWindow.SetMainMenu(menu) 
	
	self.myWindow.SetContent(container.NewHSplit(leftContentBorder, grid))
	self.myWindow.Resize(fyne.NewSize(800, 600))
	self.myWindow.CenterOnScreen()
	self.myWindow.ShowAndRun()
}

func (self *RestClient) ReloadCollections() {
	collections, err := self.DBService.ListAllCollection()
		if err != nil {
			dialog.ShowError(errors.New("Collection loading problem"), self.myWindow)
			return
		}
		self.collections = collections
}

func (self *RestClient) ExecuteRequest(urlString string, callBack func(stringRes string)) {
	jar, err := cookiejar.New(nil)		
	cookiesArray := make([]*http.Cookie, 1)
	
	if val, found := self.authDataMap["authKeyLoc"]; found && val != "Header" {
		if strings.Index(urlString, "?") > -1 {
			urlString = urlString + "&"+ self.authDataMap["authKey"]+"="+self.authDataMap["authValue"]
		} else {
			urlString = urlString + "?"+ self.authDataMap["authKey"]+"="+self.authDataMap["authValue"]
		}
	} else {
		if self.authMethod == "API Key" {
			self.dataMap[self.authDataMap["authKey"]] = self.authDataMap["authValue"]
		} else if self.authMethod == "Bearer Token" {
			self.dataMap["Authorization"] = "Bearer " + self.authDataMap["authValue"]
		} else if self.authMethod == "Basic Auth" {
			encodedString := base64.StdEncoding.EncodeToString([]byte(self.authDataMap["authKey"] + ":" + self.authDataMap["authValue"]))
			self.dataMap["Authorization"] = "Basic " + encodedString
		}
	}
	
	if (self.httpMethod == "GET") {
		self.req, err = http.NewRequest(http.MethodGet, urlString, nil)
	} else if (self.httpMethod == "HEAD") {
		self.req, err = http.NewRequest(http.MethodHead, urlString, nil)
	} else if (self.httpMethod == "POST") {
		
		if self.radioOpt1 == "x-www-form-urlencoded" {
			bodyData := url.Values{}
			if len(self.bDataMap) > 0 {
				for key, value := range self.bDataMap {
					bodyData.Add(key, value)
				}
			}
			self.req, err = http.NewRequest(http.MethodPost, urlString, strings.NewReader(bodyData.Encode()))
			self.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else if self.radioOpt1 == "form-data" {
				values := map[string]io.Reader{}
				if len(self.bDataMap) > 0 {
					for key, value := range self.bDataMap {
						values[key] = strings.NewReader(value)
					}
				}
				
				if self.selectedFile != "" {
					fileData, err := os.OpenFile(self.selectedFile, os.O_RDWR, 0755)//os.Open(self.selectedFile)
					if err != nil {
							dialog.ShowError(err, self.myWindow)
					}
					
					values["file"] = fileData
				}
				
				var dataBuffer bytes.Buffer
				multipartWriter := multipart.NewWriter(&dataBuffer)			
				defer multipartWriter.Close()
				
				for key, value := range values {
					var fw io.Writer
					if x, ok := value.(io.Closer); ok {
						defer x.Close()
					}
					if x, ok := value.(*os.File); ok {
						if fw, err = multipartWriter.CreateFormFile(key, x.Name()); err != nil {
							dialog.ShowError(err, self.myWindow)
							return
						}
					} else {
						if fw, err = multipartWriter.CreateFormField(key); err != nil {
							dialog.ShowError(err, self.myWindow)
							return
						}
					}
					if _, err = io.Copy(fw, value); err != nil {
						dialog.ShowError(err, self.myWindow)
						return
					}
				}	

				self.req, err = http.NewRequest(http.MethodPost, urlString, &dataBuffer)
				self.req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
		} else if self.radioOpt1 == "raw" {
			if self.radioOpt2 == "JSON" {
				jsonBody := bytes.NewBuffer([]byte(self.paramTextV))
				self.req, err = http.NewRequest(http.MethodPost, urlString, jsonBody)
			} else {
				self.req, err = http.NewRequest(http.MethodPost, urlString, nil)
			}
		} else if self.radioOpt1 == "binary" {
			fileData, err := os.Open(self.selectedFile)
			defer fileData.Close()
			if err != nil {
				dialog.ShowError(err, self.myWindow)
				return
			}
			
			self.req, err = http.NewRequest(http.MethodPost, urlString, fileData)
		} else {
			self.req, err = http.NewRequest(http.MethodPost, urlString, nil)
		}
	} else if (self.httpMethod == "PUT") {
		self.req, err = http.NewRequest(http.MethodPut, urlString, nil)
	} else if (self.httpMethod == "PATCH") {
		self.req, err = http.NewRequest(http.MethodPatch, urlString, nil)
	} else if (self.httpMethod == "DELETE") {
		self.req, err = http.NewRequest(http.MethodDelete, urlString, nil)
	} else if (self.httpMethod == "CONNECT") {
		self.req, err = http.NewRequest(http.MethodConnect, urlString, nil)
	} else if (self.httpMethod == "OPTIONS") {
		self.req, err = http.NewRequest(http.MethodOptions, urlString, nil)
	} else if (self.httpMethod == "TRACE") {
		self.req, err = http.NewRequest(http.MethodTrace, urlString, nil)
	} 
	
	if err != nil {
		dialog.ShowError(err, self.myWindow)
		return
	}
	
	if len(self.dataMap) > 0 {
		for key, value := range self.dataMap {
			self.req.Header.Set(key, value)
		}
	}
	
	if len(self.cDataMap) > 0 {
		self.client = &http.Client{
			Jar: jar,
		}	
		
		for key, value := range self.cDataMap {
			cookiesArray = append(cookiesArray, &http.Cookie{
				Name:   key,
				Value:  value,
				MaxAge: 300,
			})
		}				
		
		urlObj, _ := url.Parse(urlString)
		self.client.Jar.SetCookies(urlObj, cookiesArray)
		
		self.res, err = self.client.Do(self.req)
	} else {
		self.client = &http.Client{}
		self.res, err = self.client.Do(self.req)
	}
	
	if err != nil {
		dialog.ShowError(err, self.myWindow)
		return
	}

	defer self.res.Body.Close()
	
	if self.res.StatusCode == http.StatusOK {
		body, readErr := ioutil.ReadAll(self.res.Body)
		if readErr != nil {
			dialog.ShowError(readErr, self.myWindow)
			return
		}
			
		log.Println(string(body))
		callBack(string(body))
	} else {
		dialog.ShowError(errors.New("Non-OK HTTP status: " + strconv.Itoa(self.res.StatusCode)), self.myWindow)
	}
}

func main() {
	restClient := &RestClient{
		pInputURL:    "",
		httpMethod:   "",
		radioOpt1:    "",
		radioOpt2:    "",
		paramTextV:   "",
		authMethod:   "",
		authKeyLoc:   "",
		selectedFile: "",
		selColIDX: 0,
		selReqIDX: 0,
		DBService:  &services.DBService{},
	}
	
	restClient.DBService.CreateDBConnection()
	restClient.BuildUI()
}
