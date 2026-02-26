package schema

import (
	"fmt"
	"slices"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/sirupsen/logrus"
)

/*
The following pipeline.go will hold the logic for the evaluation of pipelines.
A pipeline is a series of data processing steps. Each step takes in data, processes it, and passes it to the next step.
This is useful for transforming and analyzing data in a structured way.

Pipeline acts like a DAG to process data from multiple sources and transform it stage by stage.
/*
[
{
	stageId: "stage1",
	stageType: "extract",
	dataSource: "PostgresDB",
	databaseName: "metrics_db",
	query: "SELECT user_id, activity_type, activity_timestamp FROM user_activities WHERE activity_timestamp >= NOW() - INTERVAL '1 DAY';"
	queryResultPersistTime: 0	// no caching
},
{
	stageId: "stage2",
	stageType: "foreach",
	inputStageIds: ["stage1"],
	dataSource: "PostgresDB",
	databaseName: "metrics_db",
	query: "SELECT 1 as trips, 2 as rides"
	queryResultPersistTime: 0	// no caching
},
{
	stageId: "stage3",
	stageType: "rename",
	inputStageIds: ["stage2"],
	oldColumnName: "rides",
	newColumnName: "total_rides"
},
{
	stageId: "stage3",
	stageType: "rename",
	inputStageIds: ["stage3"],
	oldColumnName: "trips",
	newColumnName: "total_trips"
},
]
]
*/
type Pipeline struct { 
	Logger *logrus.Logger
	Stages []Stage
	StageGraph map[string]Stage // adjacency list representation of the DAG
}

func (p *Pipeline) BuildPipeline(pipelineConfiguration []map[string]interface{}, dataSourceMap map[string]datasource.IDataSource) {
	// Add stages to a Graph
	p.Stages = []Stage{}
	p.StageGraph = make(map[string]Stage)
	for _, stageConfig := range pipelineConfiguration {
		stage, err := NewStage(stageConfig, dataSourceMap)
		if err != nil {
			p.Logger.Errorf("error creating stage: %v", err)
		}
		// if not in map then add it
		if _, ok :=  p.StageGraph[ stage.GetBaseStage().StageId ]; !ok {
			p.StageGraph[ stage.GetBaseStage().StageId ] = stage
		} 
	}
	// connect the stages
	for key, stageConfig := range p.StageGraph {
		inputStages := stageConfig.GetBaseStage().GetInputStages()
		for _, inputStage := range inputStages {
			st, ok := p.StageGraph[inputStage.GetBaseStage().StageId]
			if ok {
				stageConfig.SetInputStages( append( stageConfig.GetBaseStage().GetInputStages(), st) )
			}else{
				p.Logger.Errorf("input stage %s not found for stage %s", inputStage.GetBaseStage().StageId, key)
			}
		}
		
		p.Stages = append(p.Stages, stageConfig)
	}
}
func (p *Pipeline) RunPipeline() *Stage {
	// Sort pipeline stages in topological order
	// Evaluate each stage in order
	fmt.Println(p.Stages, p.StageGraph)
	sorted := p.topologicalSort()
	fmt.Println("sorted", sorted)
	var df dataframe.DataFrame;
	var err error;
	for _, stage := range sorted {
		p.Logger.Infof("Evaluating stage %s of type %s", stage.GetBaseStage().StageId, stage.GetBaseStage().StageType)
		df, err = stage.Evaluate()
		fmt.Println(df)
		if err != nil{
			p.Logger.Error(err)
		}
		stage.GetBaseStage().SetOutputDataframe(df)
	}
	fmt.Println(df, err)
	return nil
 }
func dfs(node *Stage, visited map[string]bool) []Stage {
	// if node is visited:
		// reutrn
	result := []Stage{}
	nodeAddr := fmt.Sprintf("%x", node)
	fmt.Println(nodeAddr, node, result)
	if _, ok := visited[nodeAddr]; ok {
		return []Stage{};
	}
	// add to visited
	visited[nodeAddr] = true
	
	inputs := (*node).GetBaseStage().GetInputStages()
	//for each neighbor
	for _, input := range inputs{
		//  recruse that node
		if ret := dfs(&input, visited); len(ret) > 0{
			result = append(result, ret...)
		} 
	}
	result = append(result, *node)
	return result
}

 func (p *Pipeline) topologicalSort() []Stage {
	// do topological sort of the stages in the pipeline graph
	visited := map[string]bool{}
	sorted_stages := []Stage{}
	for _, node := range p.StageGraph{
		ret := dfs(&node, visited)
		fmt.Println("ret", ret, node)
		sorted_stages = append(sorted_stages, ret...)
	}
	slices.Reverse(sorted_stages)
	return sorted_stages
}

 func NewPipeline(logger *logrus.Logger) *Pipeline {
	return &Pipeline{
		Logger: logger,
	}
}
