// ----------------------------------------------------------------------------
// Project: wx-example
// File: tenant_controller.go
// Description: Defines the Tenant controller and related authentication logic.
// Author: Your Name <you@example.com>
// Created: 2025-09-24
//
// License: MIT License
// Copyright (c) 2025 Your Name
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// ----------------------------------------------------------------------------

package controllers

import (
	"core"

	"github.com/vn-go/wx"
)

/*
All actions in Tenant require SystemUser claim
*/
type SystemUser struct {
	UserId string
}

/*
Tenant controller
*/
type Tenant struct {
	/*
		Requires SystemUser claim
	*/
	wx.Authenticate[SystemUser]
}

/*
Create a tenant if it does not exist.
@h: specifies that this method is a Web API endpoint
*/
func (t *Tenant) Create(h wx.Handler, data struct {
	Name string `json:"name" chec:"range:[5:50]"`
}) (any, error) {
	if t.Authenticate.Data == nil {
		panic("Errpr")
	}
	core.Services.TenantSvc.CreateTenant(h().Req.Context(), data.Name, data.Name)
	return data, nil
}

/*
Initialize the Tenant controller.
Sets up authentication with SystemUser claim.
*/
func (t *Tenant) New() error {
	t.Verify(func(h wx.Handler) (*SystemUser, error) {
		// get header of auth from header
		req := h().Req
		authorization := req.Header["Authorization"]

		if len(authorization) == 0 {
			return nil, wx.Errors.NewUnauthorizedError()
		}
		//s := "xknKGvzDI-sZwlXwUo2_GVsY6ce94AC3I6qNnxnOtq655tOgbcRbRnK0fs_tEb-6yz-EtjBCC0qdeDg0xu6uiw"
		user, tenant, err := core.Services.AuthSvc.Verify(req.Context(), authorization[0])
		if err != nil || user == nil || tenant != "admin" {
			//just accept at admin tenant, admin tenant is current db
			return nil, wx.Errors.NewUnauthorizedError()
		}

		return &SystemUser{
			UserId: user.UserId,
		}, nil
	})
	return nil
}
