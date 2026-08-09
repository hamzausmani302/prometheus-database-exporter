package schema

import (
	"fmt"
	"slices"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
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
	StageId         string    `yaml:"stageId"`
	StageType       StageType `yaml:"stageType"`
	InputStageIds   []string  `yaml:"inputStageIds"`
	inputStages     []Stage
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

// --- ExtractStage ---

type ExtractStage struct {
	BaseStage  `yaml:",inline"`
	Query      string `yaml:"query"`
	DataSource string `yaml:"dataSource"`
	dataSource datasource.IDataSource
}

func (ExtractStage) GetType() StageType { return StageTypeExtract }

func (e *ExtractStage) Evaluate() (dataframe.DataFrame, error) {
	if err := e.dataSource.Connect(); err != nil {
		return dataframe.DataFrame{}, fmt.Errorf("extract stage %q: connect: %w", e.StageId, err)
	}
	df := e.dataSource.GetData(datasource.SQLQuery{Query: e.Query})
	return df, nil
}
func (e *ExtractStage) GetBaseStage() *BaseStage { return &e.BaseStage }
func (e *ExtractStage) SetInputStages(s []Stage) { e.inputStages = s }

// --- ForEachStage ---
//
// For each row in the single input stage's output DataFrame, ForEachStage
// substitutes {{column}} placeholders in the query with that row's values,
// executes the query, and enriches the result with any input columns that are
// not already present. All per-row results are concatenated and returned.

type ForEachStage struct {
	BaseStage  `yaml:",inline"`
	Query      string `yaml:"query"`
	DataSource string `yaml:"dataSource"`
	dataSource datasource.IDataSource
}

func (ForEachStage) GetType() StageType { return StageTypeForeach }

func (f *ForEachStage) Evaluate() (dataframe.DataFrame, error) {
	if len(f.inputStages) == 0 {
		return dataframe.DataFrame{}, fmt.Errorf("foreach stage %q: no input stage configured", f.StageId)
	}
	if err := f.dataSource.Connect(); err != nil {
		return dataframe.DataFrame{}, fmt.Errorf("foreach stage %q: connect: %w", f.StageId, err)
	}

	inputDf := f.inputStages[0].GetBaseStage().GetOutput()
	cols := inputDf.Names()
	records := inputDf.Records() // records[0] is the header row
	if len(records) <= 1 {
		return dataframe.DataFrame{}, nil
	}

	var resultFrames []dataframe.DataFrame
	for _, row := range records[1:] {
		// Build substitution params from this input row
		params := make(map[string]string, len(cols))
		for i, col := range cols {
			params[col] = row[i]
		}

		// Substitute {{column}} placeholders in the query template
		q := f.Query
		for k, v := range params {
			q = strings.ReplaceAll(q, "{{"+k+"}}", v)
		}

		subDf := f.dataSource.GetData(datasource.SQLQuery{Query: q})
		if subDf.Nrow() == 0 {
			continue
		}

		// Propagate input-row columns that aren't already in the subquery result
		for i, colName := range cols {
			if !slices.Contains(subDf.Names(), colName) {
				vals := make([]string, subDf.Nrow())
				for k := range vals {
					vals[k] = row[i]
				}
				subDf = subDf.Mutate(series.New(vals, series.String, colName))
			}
		}
		resultFrames = append(resultFrames, subDf)
	}

	if len(resultFrames) == 0 {
		return dataframe.DataFrame{}, nil
	}
	result := resultFrames[0]
	for _, df := range resultFrames[1:] {
		result = result.RBind(df)
	}
	return result, nil
}
func (f *ForEachStage) GetBaseStage() *BaseStage { return &f.BaseStage }
func (f *ForEachStage) SetInputStages(s []Stage) { f.inputStages = s }

// --- RenameStage ---

type RenameStage struct {
	BaseStage `yaml:",inline"`
	OldName   string `yaml:"oldColumnName"`
	NewName   string `yaml:"newColumnName"`
}

func (RenameStage) GetType() StageType { return StageTypeRename }

func (rs RenameStage) Evaluate() (dataframe.DataFrame, error) {
	if len(rs.inputStages) == 0 {
		return dataframe.DataFrame{}, fmt.Errorf("rename stage %q: no input stage configured", rs.StageId)
	}
	outputDf := rs.inputStages[0].GetBaseStage().GetOutput()
	if outputDf.Nrow() == 0 {
		return dataframe.DataFrame{}, fmt.Errorf("rename stage %q: input dataframe is empty", rs.StageId)
	}
	return outputDf.Rename(rs.NewName, rs.OldName), nil
}
func (r *RenameStage) GetBaseStage() *BaseStage { return &r.BaseStage }
func (r *RenameStage) SetInputStages(s []Stage) { r.inputStages = s }

// --- Factory ---

func NewStage(stageConfig map[string]interface{}, dataSourceMap map[string]datasource.IDataSource) (Stage, error) {
	bytes, err := yaml.Marshal(stageConfig)
	if err != nil {
		return nil, err
	}

	var base BaseStage
	if err = yaml.Unmarshal(bytes, &base); err != nil {
		return nil, err
	}

	switch base.StageType {
	case StageTypeExtract:
		var s ExtractStage
		if err = yaml.Unmarshal(bytes, &s); err != nil {
			return nil, err
		}
		ds, ok := dataSourceMap[s.DataSource]
		if !ok {
			return nil, fmt.Errorf("extract stage %q: datasource %q not found", s.StageId, s.DataSource)
		}
		s.dataSource = ds
		return &s, nil

	case StageTypeForeach:
		var s ForEachStage
		if err = yaml.Unmarshal(bytes, &s); err != nil {
			return nil, err
		}
		ds, ok := dataSourceMap[s.DataSource]
		if !ok {
			return nil, fmt.Errorf("foreach stage %q: datasource %q not found", s.StageId, s.DataSource)
		}
		s.dataSource = ds
		return &s, nil

	case StageTypeRename:
		var s RenameStage
		if err = yaml.Unmarshal(bytes, &s); err != nil {
			return nil, err
		}
		return &s, nil

	default:
		return nil, fmt.Errorf("unknown stage type %q", base.StageType)
	}
}
