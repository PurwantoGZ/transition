package transition

import (
	"fmt"
	"github.com/go-pg/pg/v10/orm"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/audited"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
)

// StateChangeLog a model that used to keep state change logs
type StateChangeLog struct {
	gorm.Model
	ReferTable string
	ReferID    string
	From       string
	To         string
	Note       string `sql:"size:1024"`
	audited.AuditedModel
}

// GenerateReferenceKey generate reference key used for change log
func GenerateReferenceKey(model interface{}, db orm.DB) string {
	var (
		scope         = orm.GetTable(reflect.TypeOf(model).Elem())
		primaryValues []string
	)

	for _, field := range scope.PKs {
		primaryValues = append(primaryValues, fmt.Sprint(field.SQLName))
	}

	return strings.Join(primaryValues, "::")
}

// GetStateChangeLogs get state change logs
func GetStateChangeLogs(model interface{}, db orm.DB) []StateChangeLog {
	var (
		changelogs []StateChangeLog
		scope      = orm.GetTable(reflect.TypeOf(model).Elem())
	)

	err := db.Model(&changelogs).Where("refer_table = ? AND refer_id = ?", scope.ModelName, GenerateReferenceKey(model, db)).Select()
	if err != nil {
		panic(err)
	}
	return changelogs
}

// GetLastStateChange gets last state change
func GetLastStateChange(model interface{}, db orm.DB) *StateChangeLog {
	var (
		changelog StateChangeLog
		scope      = orm.GetTable(reflect.TypeOf(model).Elem())
	)

	err := db.Model(&changelog).Where("refer_table = ? AND refer_id = ?", scope.ModelName, GenerateReferenceKey(model, db)).Last()
	if err != nil {
		panic(err)
	}
	if changelog.To == "" {
		return nil
	}
	return &changelog
}

// ConfigureQorResource used to configure transition for qor admin
func (stageChangeLog *StateChangeLog) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		if res.Permission == nil {
			res.Permission = roles.Deny(roles.Update, roles.Anyone).Deny(roles.Create, roles.Anyone)
		} else {
			res.Permission = res.Permission.Deny(roles.Update, roles.Anyone).Deny(roles.Create, roles.Anyone)
		}
	}
}
