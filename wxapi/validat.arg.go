/*
This file give a exmaple how to validate the input arguments in wx controller.
*/
package main

import "github.com/vn-go/wx"

type UserData struct {
	Username string `json:"username"`
}
type MyController struct {
	Name string
	//Emulator user store
	Userstore map[string]UserData
	// Emulator user wallet
	UserWallet map[string]float64
}

func (c *MyController) New() {
	c.Userstore = map[string]UserData{
		"admin": UserData{Username: "admin"},
		"user1": UserData{Username: "user1"},
	}
	c.UserWallet = map[string]float64{
		"admin": 1000000,
		"user1": 500000,
	}
}

/*
	Note: The func in controller do not have wx.Handler argument is a logic function in controller.
	This function is not hanlded by wx.Router even if this function is exported.
*/
func (c *MyController) FindUser( //<-- even function name is exported, it is not hanlded by wx.Router
	name string,
) *UserData {
	// emulate find user by name
	if u, ok := c.Userstore[name]; ok {
		return &u
	} else {
		return nil
	}
}
func (c *MyController) CheckWallet( //<-- even function name is exported, it is not hanlded by wx.Router
	name string,
	amount float64,
) bool {
	// emulate find user by name
	if cashInWallet, ok := c.UserWallet[name]; ok {
		return cashInWallet >= amount
	} else {
		return false
	}
}

/*
	This is an example of validating the input arguments in wx controller.
	Note: all all invalid check will return 400 Bad Request error before handler called.
	wx will check automatically the input arguments before calling the handler.
*/
func (c *MyController) AutoValidateArgApi(
	h *wx.Handler, //<-- this is API route handler
	data struct {
		Name      string  `json:"name" check:"range:[3:20]"`                                           //<-- leng of name should be between 3 and 20
		Email     string  `json:"email" check:"regex:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"` //<-- email should be valid format
		Age       int     `json:"age" check:"range:[18:120]"`                                          //<-- age should be between 18 and 120
		Price     float64 `json:"price" check:"range:[0.12:300]"`
		TotalItem int     `json:"total_item" check:"range:[1:100]"` //<-- price should be between 0.12 usd and 300usd
		TotalCash float64 `json:"-"`                                //<-- not show in api doc
	},
) (any, error) {
	// remmeber that wx will check automatically the input arguments before calling the handler.
	user := c.FindUser(data.Name)
	if user == nil {
		// use wx.Errors return proper error message to client
		return nil, wx.Errors.NewForbidenError("User not found")
	}
	data.TotalCash = data.Price * float64(data.TotalItem)
	if !c.CheckWallet(data.Name, data.TotalCash) {
		return nil, wx.Errors.NewHttpError(402, "Payment amount is too large")
	}
	return data, nil
}
