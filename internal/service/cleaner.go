package service

import (
	"go-openstackclient/consts"
	"go-openstackclient/internal"
	"go-openstackclient/utils"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"
)


type Cleaner struct {
	adminManager                   *Controller
	runners                        []*ProjectRunner
	wg                             sync.WaitGroup
}

func initProjectRunners(controller *Controller, projects []string) []*ProjectRunner {
	toDeleteProjects := checkProjectExist(controller, projects)
	runners := make([]*ProjectRunner, 0)
	for _, projectName := range toDeleteProjects {
		projectRunner := NewProjectRunner(projectName)
		runners = append(runners, projectRunner)
	}
	return runners
}

func NewCleaner(projects []string) *Cleaner {
	adminManager := &Controller{
		projectName: consts.ADMIN,
	}
	return &Cleaner{
		adminManager: adminManager,
		runners: initProjectRunners(adminManager, projects),
	}
}

func checkProjectExist(controller *Controller, projects []string) []string {
	toDeleteProjects := make([]string, 0)
	for _, projectName := range projects {
		if controller.GetProjectId(projectName) == "" {
		    log.Printf("@@@@@@@@@@@@@@@Project %s not exist, not to delete resources\n", projectName)
		} else {
			toDeleteProjects = append(toDeleteProjects, projectName)
		}
	}
	return toDeleteProjects
}

func (c *Cleaner) Run() {
	for _, runner := range c.runners {
    	c.wg.Add(1)
    	go runner.Run(&c.wg)
	}
	c.wg.Wait()
    c.report()
}

func (c *Cleaner) report() {
	for _, runner := range c.runners {
		runner.makeReport()
	}
}

type ProjectRunner struct {
	projectName        string
	manager            *Controller
	depNodes           map[string]*Node
	completedChannel   chan struct{}
}

func NewProjectRunner(projectName string) *ProjectRunner {
	m := NewController(projectName)
	depNodes := InitNodes()
	return &ProjectRunner{
		projectName: projectName,
		manager: m,
		depNodes: depNodes,
		completedChannel: make(chan struct{}, len(depNodes)),
	}
}

func (p *ProjectRunner) getMethodName(resourceType string) string {
	resourceType = utils.Pluralize(resourceType)
	var res = "Delete"
	for _, s := range strings.Split(resourceType, "_") {
		res += strings.Title(s)
	}
	return res
}

func (p *ProjectRunner) callNode(node *Node) {
	methodName := p.getMethodName(node.resourceType)
	defer func() {
		if err := recover(); err != nil {
			log.Println("call error occur", err)
		}
		for _, dep := range node.dependencies {
			p.depNodes[dep.resourceType].monitorDeleteChannel <- struct{}{}
			log.Printf("%s %s completed, notify %s", node.resourceType, methodName, dep.resourceType)
		}
		log.Printf("%s call completed", node.resourceType)
		p.completedChannel <- struct{}{}
	}()

	if len(node.monitorDeleteChannel) != cap(node.monitorDeleteChannel) {
		for len(node.monitorDeleteChannel) != cap(node.monitorDeleteChannel) {}
	}
	log.Printf("Cleaning %s is in progress", node.resourceType)
	reflect.ValueOf(p.manager).MethodByName(methodName).Call([]reflect.Value{})
}

func (p *ProjectRunner) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for resourceType, _ := range p.depNodes {
		go p.callNode(p.depNodes[resourceType])
	}

	for len(p.completedChannel) != cap(p.completedChannel) {
		time.Sleep(2 * time.Second)
		log.Println("waiting for completed...")
	}
	//p.manager.DeleteUserByName(p.projectName)
	//p.manager.DeleteProjectByName(p.projectName)
	log.Printf("@@@@@@@@@@@@@@@Clean project %s completed", p.projectName)
}

func mergeMaps(map1, map2 map[string]chan internal.Output) map[string]chan internal.Output {
	mergedMap := make(map[string]chan internal.Output)

	for key, value := range map1 {
		mergedMap[key] = value
	}

	for key, value := range map2 {
		mergedMap[key] = value
	}

	return mergedMap
}

func (p *ProjectRunner) makeReport() {
	outputs := p.manager.DeleteChannels
	reporters := make(map[string]reporter)
	for resourceType, ch := range outputs {
		totals := len(ch)
		close(ch)
		failed, succeed := make([]Output, 0), make([]map[string]string, 0)
		for output := range ch {
			if !output.Success {
				failed = append(failed, output)
			} else {
				succeed = append(succeed, output.ParametersMap)
			}
		}
		r := reporter{
			resourceType: resourceType,
			totals: totals,
			failed: failed,
			succeed: succeed,
		}
		reporters[resourceType] = r
	}
	log.Printf("Project %s reported:***********************************************\n", p.projectName)
	for _, resourceType := range OrderResources {
		output := reporters[resourceType]
        log.Printf("Resource %-*s-----> totals %d, succeed %s, failed %+v\n",
        	25, resourceType, output.totals, output.succeed, output.failed)
	}
}

type reporter struct {
	resourceType            string
	totals                  int
	failed                  []Output
	succeed                 []map[string]string
}

