// +build ignore

package mediastorage

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
	// SystemMediaStorageConfigurationMediaStorage => Media Storage.
	// Path: system/media_storage_configuration/media_storage
	// SourceModel: Otnegam\MediaStorage\Model\Config\Source\Storage\Media\Storage
	SystemMediaStorageConfigurationMediaStorage model.Str

	// SystemMediaStorageConfigurationMediaDatabase => Select Media Database.
	// Path: system/media_storage_configuration/media_database
	// BackendModel: Otnegam\MediaStorage\Model\Config\Backend\Storage\Media\Database
	// SourceModel: Otnegam\MediaStorage\Model\Config\Source\Storage\Media\Database
	SystemMediaStorageConfigurationMediaDatabase model.Str

	// SystemMediaStorageConfigurationSynchronize => .
	// After selecting a new media storage location, press the Synchronize button
	// to transfer all media to that location. Media will not be available in the
	// new location until the synchronization process is complete.
	// Path: system/media_storage_configuration/synchronize
	SystemMediaStorageConfigurationSynchronize model.Str

	// SystemMediaStorageConfigurationConfigurationUpdateTime => Environment Update Time.
	// Path: system/media_storage_configuration/configuration_update_time
	SystemMediaStorageConfigurationConfigurationUpdateTime model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemMediaStorageConfigurationMediaStorage = model.NewStr(`system/media_storage_configuration/media_storage`, model.WithConfigStructure(cfgStruct))
	pp.SystemMediaStorageConfigurationMediaDatabase = model.NewStr(`system/media_storage_configuration/media_database`, model.WithConfigStructure(cfgStruct))
	pp.SystemMediaStorageConfigurationSynchronize = model.NewStr(`system/media_storage_configuration/synchronize`, model.WithConfigStructure(cfgStruct))
	pp.SystemMediaStorageConfigurationConfigurationUpdateTime = model.NewStr(`system/media_storage_configuration/configuration_update_time`, model.WithConfigStructure(cfgStruct))

	return pp
}
