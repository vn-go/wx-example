package core

import (
	"context"
	"core/models"
	"fmt"
	"strings"
	"time"

	"github.com/vn-go/dx"
)

type MenuItem struct {
	Title    string     `json:"title"`
	Icon     string     `json:"icon"`
	Link     string     `json:"link,omitempty"`     // chỉ có khi là menu con
	Children []MenuItem `json:"children,omitempty"` // lồng nhau nhiều cấp
}
type systemService struct {
	tenantSvc tenantService
}

func (s *systemService) SyncMenu(context context.Context, user *UserClaims, menuData []MenuItem) error {
	type appMenuKeyInfo struct {
		Id       uint64
		ParentId *uint64
		IdPaths  string
	}
	var err error
	if menuData == nil {
		return nil
	}
	// convert menuData to models.AppMenu
	//var appMenu []models.AppMenu
	err = s.convertMenuToAppMenu(menuData, user, func(db *dx.DB, appMenu *models.AppMenu) error {

		err = db.Insert(appMenu)
		if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
			if dbErr.IsDuplicateEntryError() {
				item, err := dx.QueryItem[appMenuKeyInfo](db, "appMenu(id Id,parentId,IdPaths),where(viewPath=?)", appMenu.ViewPath)
				if err != nil {
					return err
				}
				if appMenu != nil {
					appMenu.ID = item.Id
					appMenu.ParentId = item.ParentId
					if appMenu.IdPaths == "" {
						appMenu.IdPaths = fmt.Sprintf(".%d.", appMenu.ID)
					}
					if appMenu.ParentId != nil {
						parentItem, err := dx.QueryItem[appMenuKeyInfo](db, "appMenu(id Id,parentId,IdPaths),where(id=?)", appMenu.ParentId)
						if err != nil {
							return err
						}
						appMenu.IdPaths = fmt.Sprintf("%s%d.", parentItem.IdPaths, appMenu.ID)

					}
					db.Update(&appMenu)

				} else {
					return fmt.Errorf("cannot find menu with viewPath %s", appMenu.ViewPath)
				}

			} else {
				return err
			}
		} else {
			appMenu.IdPaths = fmt.Sprintf(".%d.", appMenu.ID)
			if appMenu.ParentId != nil {
				parentItem, err := dx.QueryItem[appMenuKeyInfo](db, "appMenu(id Id,parentId,IdPaths),where(id=?)", appMenu.ParentId)
				if err != nil {
					return err
				}
				appMenu.IdPaths = fmt.Sprintf("%s%d.", parentItem.IdPaths, appMenu.ID)
				// appMenu.ViewPath = fmt.Sprintf(".%d%s", appMenu.ID, parentItem.IdPaths)

			}
			r := db.Update(appMenu)
			if r.Error != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// save to database

	return nil
}

// func (s *systemService) saveAppMenu(appMenu []models.AppMenu, user *UserClaims) error {
// 	db, err := s.tenantSvc.GetTenant(user.Tenant)
// 	if err != nil {
// 		return err
// 	}
// 	for _, menu := range appMenu {
// 		err := db.Insert(&menu)
// 		if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
// 			if dbErr.IsDuplicateEntryError() {
// 				dx.QueryItems[uint64](db, "appMenu(id),where(viewPath=?)", menu.ViewPath)
// 			}
// 		}
// 	}
// 	panic("unimplemented")
// }

func (s *systemService) convertMenuToAppMenu(menuData []MenuItem, user *UserClaims, onSave func(db *dx.DB, appMenu *models.AppMenu) error) error {
	db, err := s.tenantSvc.GetTenant(user.Tenant)
	if err != nil {
		return err
	}
	// Hàm đệ quy để convert menu
	var convertRecursive func(menuItems []MenuItem, parentID *uint64, createdBy string) error
	convertRecursive = func(menuItems []MenuItem, parentID *uint64, createdBy string) error {
		// var result []models.AppMenu
		currentTime := time.Now()

		for _, item := range menuItems {
			appMenu := models.AppMenu{

				ParentId:  parentID,
				Title:     item.Title,
				Icon:      item.Icon,
				ViewPath:  strings.ToLower(item.Link),
				CreatedBy: createdBy,
				CreatedOn: currentTime,
				UpdatedBy: createdBy,
				UpdatedOn: nil,
			}
			err := onSave(db, &appMenu)
			if err != nil {
				return err
			}

			// Add parent menu to result
			// Thêm menu cha vào kết quả
			// result = append(result, appMenu)
			// resole sub menu if existing
			// Xử lý menu con nếu có
			if len(item.Children) > 0 {
				// Creae a parentID for sub menu
				// Tạo một parentID giả định cho menu con
				// Trong thực tế, bạn cần biết ID thực của menu cha
				tempParentID := appMenu.ID
				err := convertRecursive(item.Children, &tempParentID, createdBy)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	return convertRecursive(menuData, nil, user.Username)
}

func newSystemService(tenantSvc tenantService) *systemService {
	return &systemService{
		tenantSvc: tenantSvc,
	}
}
