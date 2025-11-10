package core

type tamperService struct {
	tenantSvc *tenantService
	cacheSvc  cacheService
}

func newTamperService(tenantSvc tenantService, cacheSvc cacheService) *tamperService {
	return &tamperService{
		tenantSvc: &tenantSvc,
		cacheSvc:  cacheSvc,
	}
}
