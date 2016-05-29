package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/maleck13/templator/model"
	k8 "k8s.io/kubernetes/pkg/api/v1"
)

const TEMPLATES_FILE_LOC = "./.templates.json"

//may add support for a db if wanted to use as a lib but focus on cli for now
type TemplateService struct {
	DataType string
}

func NewTemplateService(dataType string) *TemplateService {
	return &TemplateService{DataType: dataType}
}

func (ts *TemplateService) GetTemplate(name string) (*model.ApplicationTemplate, error) {
	if ts.DataType == "local" {
		templates, err := loadDataFromFile(TEMPLATES_FILE_LOC)
		if err != nil {
			return nil, err
		}
		return templates[name], nil
	}
	return nil, errors.New("unsupported data type")
}

func (ts *TemplateService) ListTemplates() (map[string]*model.ApplicationTemplate, error) {
	return loadDataFromFile(TEMPLATES_FILE_LOC)
}

func (ts *TemplateService) SaveTemplate(name string, tempModel *model.ApplicationTemplate) error {
	data, err := loadDataFromFile(TEMPLATES_FILE_LOC)
	if err != nil {
		return err
	}
	data[name] = tempModel
	return saveDataToFile(TEMPLATES_FILE_LOC, data)

}

func (ts *TemplateService) DeleteTemplate(name string) error {
	data, err := loadDataFromFile(TEMPLATES_FILE_LOC)
	if err != nil {
		return err
	}
	delete(data, name)
	return saveDataToFile(TEMPLATES_FILE_LOC, data)
}

func (ts *TemplateService) SaveDeployment(tempName, depName string, tempModel *model.OSTDeploymentConfig) error {
	data, err := loadDataFromFile(TEMPLATES_FILE_LOC)
	if appTemp, ok := data[tempName]; ok {
		if nil == appTemp.DeploymentConfigs {
			appTemp.DeploymentConfigs = make(map[string]*model.OSTDeploymentConfig)
		}
		appTemp.DeploymentConfigs[depName] = tempModel
		data[tempName] = appTemp
	}
	if err != nil {
		return err
	}
	return saveDataToFile(TEMPLATES_FILE_LOC, data)

}

func (ts *TemplateService) SaveService(tempName, depName string, tempModel *k8.Service) error {
	data, err := loadDataFromFile(TEMPLATES_FILE_LOC)
	if err != nil {
		return err
	}
	if appTemp, ok := data[tempName]; ok {
		appTemp.Services[depName] = tempModel
		data[tempName] = appTemp
	}

	return saveDataToFile(TEMPLATES_FILE_LOC, data)

}

func loadDataFromFile(location string) (map[string]*model.ApplicationTemplate, error) {
	reader, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	data := make(map[string]*model.ApplicationTemplate)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil

}

func saveDataToFile(location string, data map[string]*model.ApplicationTemplate) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(location, content, 0644)
}
