package schema

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/sirupsen/logrus"
)

/*
Pipeline executes a directed acyclic graph (DAG) of stages. Each stage receives
the output DataFrame of its declared input stages and produces a new DataFrame.
The final stage's output is returned as the pipeline result.

YAML pipeline config example:

	pipeline:
	  - stageId: wells
	    stageType: extract
	    dataSource: id3_dashboard
	    query: "SELECT id_t_well, well_name, well_id FROM t_well WHERE is_live"

	  - stageId: gaps
	    stageType: foreach
	    dataSource: id3_dashboard
	    inputStageIds: [wells]
	    query: "SELECT sum(gap) AS total_gap FROM t_mudlog WHERE id_t_well = {{id_t_well}}"

Stage types:
  - extract  — runs a plain SQL query; no inputs required
  - foreach  — for each row of the single input stage, substitutes {{column}}
    placeholders in the query, runs it, and merges results
  - rename   — renames a column in the single input stage's DataFrame
*/
type Pipeline struct {
	Logger     *logrus.Logger
	Stages     []Stage
	StageGraph map[string]Stage
}

// BuildPipeline constructs the stage DAG from raw YAML config maps.
// It performs two passes: first it creates all Stage objects and registers
// them by stageId, then it resolves inputStageIds to actual Stage references.
func (p *Pipeline) BuildPipeline(pipelineConfig []map[string]interface{}, dataSourceMap map[string]datasource.IDataSource) {
	p.Stages = []Stage{}
	p.StageGraph = make(map[string]Stage)

	// Pass 1: create and register all stages
	for _, cfg := range pipelineConfig {
		stage, err := NewStage(cfg, dataSourceMap)
		if err != nil {
			p.Logger.Errorf("pipeline: error creating stage: %v", err)
			continue
		}
		id := stage.GetBaseStage().StageId
		if _, exists := p.StageGraph[id]; exists {
			p.Logger.Warnf("pipeline: duplicate stageId %q — skipping", id)
			continue
		}
		p.StageGraph[id] = stage
	}

	// Pass 2: resolve inputStageIds → actual Stage references
	for _, stage := range p.StageGraph {
		var resolved []Stage
		for _, depId := range stage.GetBaseStage().InputStageIds {
			dep, ok := p.StageGraph[depId]
			if !ok {
				p.Logger.Errorf("pipeline: stage %q references unknown input stage %q", stage.GetBaseStage().StageId, depId)
				continue
			}
			resolved = append(resolved, dep)
		}
		stage.SetInputStages(resolved)
		p.Stages = append(p.Stages, stage)
	}
}

// RunPipeline executes all stages in topological order and returns the last
// stage's output DataFrame. Each stage receives the outputs of its input
// stages via GetBaseStage().GetOutput().
func (p *Pipeline) RunPipeline() (dataframe.DataFrame, error) {
	sorted := p.topologicalSort()
	p.Logger.Debugf("pipeline: running %d stages", len(sorted))

	var last dataframe.DataFrame
	for _, stage := range sorted {
		id := stage.GetBaseStage().StageId
		p.Logger.Infof("pipeline: evaluating stage %q (%s)", id, stage.GetType())
		df, err := stage.Evaluate()
		if err != nil {
			return dataframe.DataFrame{}, fmt.Errorf("pipeline: stage %q failed: %w", id, err)
		}
		stage.GetBaseStage().SetOutputDataframe(df)
		p.Logger.Debugf("pipeline: stage %q produced %d rows", id, df.Nrow())
		last = df
	}
	return last, nil
}

// topologicalSort returns stages in dependency-first order using DFS post-order
// traversal. Dependencies are always visited before the stages that depend on them.
func (p *Pipeline) topologicalSort() []Stage {
	visited := map[string]bool{}
	var sorted []Stage
	for _, node := range p.StageGraph {
		sorted = append(sorted, dfs(&node, visited)...)
	}
	return sorted
}

// dfs performs a depth-first post-order traversal: all inputs of a node are
// appended before the node itself, guaranteeing correct execution order.
func dfs(node *Stage, visited map[string]bool) []Stage {
	id := (*node).GetBaseStage().StageId
	if visited[id] {
		return nil
	}
	visited[id] = true

	var result []Stage
	for _, input := range (*node).GetBaseStage().GetInputStages() {
		result = append(result, dfs(&input, visited)...)
	}
	result = append(result, *node)
	return result
}

func NewPipeline(logger *logrus.Logger) *Pipeline {
	return &Pipeline{Logger: logger}
}
