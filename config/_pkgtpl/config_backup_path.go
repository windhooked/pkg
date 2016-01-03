// +build ignore

package backup

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// SystemBackupEnabled => Enable Scheduled Backup.
	// Path: system/backup/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SystemBackupEnabled model.Bool

	// SystemBackupType => Backup Type.
	// Path: system/backup/type
	// SourceModel: Otnegam\Backup\Model\Config\Source\Type
	SystemBackupType model.Str

	// SystemBackupTime => Start Time.
	// Path: system/backup/time
	SystemBackupTime model.Str

	// SystemBackupFrequency => Frequency.
	// Path: system/backup/frequency
	// BackendModel: Otnegam\Backup\Model\Config\Backend\Cron
	// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
	SystemBackupFrequency model.Str

	// SystemBackupMaintenance => Maintenance Mode.
	// Please put your store into maintenance mode during backup.
	// Path: system/backup/maintenance
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SystemBackupMaintenance model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemBackupEnabled = model.NewBool(`system/backup/enabled`, model.WithConfigStructure(cfgStruct))
	pp.SystemBackupType = model.NewStr(`system/backup/type`, model.WithConfigStructure(cfgStruct))
	pp.SystemBackupTime = model.NewStr(`system/backup/time`, model.WithConfigStructure(cfgStruct))
	pp.SystemBackupFrequency = model.NewStr(`system/backup/frequency`, model.WithConfigStructure(cfgStruct))
	pp.SystemBackupMaintenance = model.NewBool(`system/backup/maintenance`, model.WithConfigStructure(cfgStruct))

	return pp
}
