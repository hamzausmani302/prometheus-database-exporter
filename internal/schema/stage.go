package schema

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"gopkg.in/yaml.v3"
)

type StageType string

const (
	StageTypeExtract StageType = "extract"
	StageTypeForeach StageType = "foreach"
	StageTypeRename  StageType = "rename"
)

type Stage interface {
	Evaluate() (dataframe.DataFrame, error)
	GetType() StageType
	GetBaseStage() *BaseStage
	SetInputStages(stages []Stage)
}

type BaseStage struct {
	StageId      string    `yaml:"stageId"`
	StageType    StageType `yaml:"stageType"`
	inputStageIds []string    `yaml:"inputStageIds"`
	inputStages  []Stage
	outputDataframe dataframe.DataFrame
}
func (b *BaseStage) SetOutputDataframe(df dataframe.DataFrame) {
	b.outputDataframe = df
}
func (b *BaseStage) GetOutput() dataframe.DataFrame {
	return b.outputDataframe
}
func (b *BaseStage) GetInputStages() []Stage {
	return b.inputStages
}

type ExtractStage struct {
	BaseStage  `yaml:",inline"`
	Query      string `yaml:"query"`
	DataSource string `yaml:"dataSource"`
	dataSource datasource.IDataSource
}
func(ExtractStage) GetType() StageType {return StageTypeExtract}
func(e ExtractStage) Evaluate() (dataframe.DataFrame, error) {
	// read the data from datastore and return it for the next stage
	df := e.dataSource.GetData(datasource.SQLQuery{
		Query: e.Query,
	})
	return df, nil
}
func (e *ExtractStage) GetBaseStage() *BaseStage {
	return &e.BaseStage
}
func (e *ExtractStage) SetInputStages(stages []Stage) {
	e.inputStages = stages
}

type ForEachStage struct {
	BaseStage  `yaml:",inline"`
	Query      string `yaml:"query"`
	DataSource string `yaml:"dataSource"`
	dataSource datasource.IDataSource
}
func(ForEachStage) GetType() StageType {return StageTypeForeach}
func (ForEachStage) Evaluate() (dataframe.DataFrame, error) {
	// ForEachStage is not yet implemented
	return dataframe.DataFrame{}, nil
}
func (f *ForEachStage) GetBaseStage() *BaseStage {
	return &f.BaseStage
}
func (f *ForEachStage) SetInputStages(stages []Stage) {
	f.inputStages = stages
}

type RenameStage struct {
	BaseStage `yaml:",inline"`
	OldName   string `yaml:"oldColumnName"`
	NewName   string `yaml:"newColumnName"`
}

func(RenameStage) GetType() StageType {return StageTypeRename}
func(rs RenameStage) Evaluate() (dataframe.DataFrame, error) {
	if len(rs.inputStages) == 0 {
		return dataframe.DataFrame{}, fmt.Errorf("input stage is nil for rename stage %s", rs.StageId)
	}
	output_df := rs.inputStages[0].GetBaseStage().GetOutput()
	if output_df.Nrow() == 0 {
		return dataframe.DataFrame{}, fmt.Errorf("input dataframe is empty for rename stage %s", rs.StageId)
	}
	output_df = output_df.Rename( rs.NewName, rs.OldName)
	return output_df, nil
}
func (r *RenameStage) GetBaseStage() *BaseStage {
	return &r.BaseStage
}
func (r *RenameStage) SetInputStages(stages []Stage) {
	r.inputStages = stages
}



func  NewStage(stageConfig map[string]interface{}, dataSourceMap map[string]datasource.IDataSource) (Stage, error) {
	var stage Stage
	var baseStage  BaseStage;
	bytes , err := yaml.Marshal(stageConfig)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &baseStage)
	if err != nil {
		return nil, err
	}
	if (baseStage.StageType == StageTypeExtract) {
		var extractStage ExtractStage
		err = yaml.Unmarshal(bytes, &extractStage)
		extractStage.dataSource = dataSourceMap[extractStage.DataSource]
		stage = &extractStage
	}else if( baseStage.StageType == StageTypeForeach){
		var forEachStage ForEachStage
		err = yaml.Unmarshal(bytes, &forEachStage)
		forEachStage.dataSource = dataSourceMap[forEachStage.DataSource]
		stage = &forEachStage
	}else if( baseStage.StageType == StageTypeRename){
		var renameStage RenameStage
		err = yaml.Unmarshal(bytes, &renameStage)
		stage = &renameStage
	}
	if err != nil {
		return nil, err
	}
	return stage, nil
}
