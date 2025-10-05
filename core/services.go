package core

import (
	"fmt"

	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

type Container struct {
	bx.Container[Container]
	userService  bx.Depen[userService, Container]
	cacheService bx.Depen[cacheService, Container]
	userRepo     bx.Depen[userRepo, Container]
	pwdSvc       bx.OnlyOnce[passwordService]

	msg bx.OnlyOnce[bx.MessageBus]

	db bx.OnlyOnce[*dx.DB]

	config bx.OnlyOnce[*configInfo] //<-- OnlyOnce 1 lan duy nhat (singleton)
}

/*
Khoi tao ban dau cho container, moi conatiner bat buoc phai co ham New, dung de setup cac Resolver
The initialization for the container, each container must have a New function, used to set up the Resolvers.
*/
func (c *Container) New() {
	c.userRepo.Resolve(func() (userRepo, error) {

		ret := &userRepoSql{}

		return ret, nil

	})
	/*
	 setting up user service
	*/
	c.userService.Resolve(func() (userService, error) {
		var err error
		ret := &userServiceSql{}
		ret.userRepo, err = c.userRepo.Get()
		if err != nil {
			return nil, err
		}
		ret.pwdSvc, err = c.pwdSvc.Get()
		if err != nil {
			return nil, err
		}
		ret.cache, err = c.cacheService.Get()
		if err != nil {
			if err != nil {
				return nil, err
			}
		}

		return ret, nil
	})
	c.pwdSvc.Resolve(func() (passwordService, error) {
		cfg, err := c.config.Get()
		if err != nil {
			return nil, err
		}
		if cfg.Jwt.HashPasswordType == "bcypt" {
			return &bcryptPasswordService{}, nil
		}
		if cfg.Jwt.HashPasswordType == "argon" {
			return &argon2PasswordService{}, nil
		}
		return nil, fmt.Errorf("%s was not suppport", cfg.Jwt.HashPasswordType)

	})

	c.cacheService.Resolve(func() (cacheService, error) {
		ret := &cacheServiceImpl{
			cacher: bx.Cacher.NewInMemoryCache(),
		}

		return ret, nil
	})

}
func NewContainer(configPath string) *Container {
	c := (&Container{}).Get()
	c.config.Resolve(func() (*configInfo, error) {

		return loadConfig(configPath)

	})
	return c
}

// type DbService struct {
// 	Db *dx.DB
// }

// var DbSvc = (&bx.Wire[servicesTypes]{}).WireThenGet(
// 	newConfig,
// 	newDB,
// )

type servicesTypes struct {
	Config    *configInfo
	Users     userService
	UserRepo  userRepo
	Cache     cacheService
	Db        *dx.DB
	TenantSvc tenantService
	AuthSvc   serviceAuth
	Log       logService
	Encrypt   encryptService
	Broker    bx.Broker
	JWTSvc    jwtService
	RABCSvc   rabcService
	DataSvc   dataService
	//Pwd passwordService
	//config configInfo
}

var Services = (&bx.Wire[servicesTypes]{}).WireThenGet(
	newConfig,
	newDB,
	newJwtService,
	newTenantService,
	newServiceAuth,

	newUserRepoSql,
	newUserServiceSql,
	//newCacheServiceImpl,
	newCacheService,
	newpasswordService,
	newLogLumberjack,
	newEncryptServiceImpl,
	newBrokerService,
	newValidatorService,
	NewRabcService,
	newDataService,
)
