package core

import (
	"core/services/config"
	"sync"
	"time"

	"github.com/vn-go/dx"
)

var globalDb *dx.DB
var newDbOnce sync.Once

func newDb(cfgSvc *config.ConfigService) *dx.DB {
	db, err := dx.Open(cfgSvc.Get().Database.Driver, cfgSvc.Get().Database.DSN)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(80)

	db.SetMaxIdleConns(20)
	//db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(5 * time.Minute)
	return db
}

//mklink /D "C:\ProgramData\Microsoft\VisualStudio\Packages" "D:\VSPackages"
// C:\ProgramData (x86)\Microsoft VisualStudio\Shared

//robocopy "C:\Users\MSI CYBORG\AppData\Local\Microsoft\VisualStudio" "D:\Vs2026\VSAppData" /E /MOVE
//mklink /D "C:\Users\MSI CYBORG\AppData\Local\Microsoft\VisualStudio" "D:\Vs2026\VSAppData"

//robocopy "C:\Program Files (x86)\Microsoft Visual Studio\Shared" "D:\Vs2026\VSShared" /E /MOVE
//mklink /D "C:\Program Files (x86)\Microsoft Visual Studio\Shared" "D:\Vs2026\VSShared"

/*
dotnet new webapi -n app1 --framework net10.0
cd app1
dotnet add package Swashbuckle.AspNetCore --version 6.9.0  # Stable cho .NET 10
dotnet restore
*/
