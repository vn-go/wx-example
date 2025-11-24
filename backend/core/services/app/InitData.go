package app

func (app *AppService) InitData() {
	cfg := app.cfgSvc.Get()
	if cfg.Tenant.IsMulti {
		app.InitMultiTenant()

		app.createTenantConfigAccount()
	} else {
		app.InitSingleTenant()
		app.createDefaultAccount()

	}

}

func (app *AppService) createTenantConfigAccount() {
	panic("unimplemented")
}

func (app *AppService) InitMultiTenant() {
	panic("unimplemented")
}
