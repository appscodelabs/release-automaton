package main

import (
	"encoding/json"
	"fmt"
	"github.com/appscodelabs/release-automaton/lib"
	"sort"
)

type Project struct {
	Name string `json:"name"`
	Versions []string `json:"versions"`
}

type Catalog struct {
	Projects []Project `json:"projects"`
}

func (c *Catalog) Sort() {
	sort.Slice(c.Projects, func(i, j int) bool { return c.Projects[i].Name < c.Projects[j].Name })
	var err error
	for i, project := range c.Projects {
		c.Projects[i].Versions, err = lib.SortVersions(project.Versions)
		if err != nil {
			panic(err)
		}
	}
}

func CreateCatalogData() Catalog {
	return Catalog{
		Projects: []Project{
			{
				Name: "postgres",
				Versions: []string{
					"9.6",
					"10.2",
					"10.6",
					"11.1",
					"11.2",
				},
			},
			{
				Name: "mongodb",
				Versions: []string{
					"3.4.17",
					"3.4.22",
					"3.6.8",
					"3.6.13",
					"4.0.3",
					"4.0.5",
					"4.0.11",
					"4.1.4",
					"4.1.7",
					"4.1.13",
					"4.2.3",
				},
			},
			{
				Name: "elasticsearch",
				Versions: []string{
					"5.6.4",
					"6.2.4",
					"6.3.0",
					"6.4.0",
					"6.5.3",
					"6.8.0",
					"7.2.0",
					"7.3.2",
				},
			},
			{
				Name: "mysql",
				Versions: []string{
					"5.7.25",
					"8.0.3",
					"8.0.14",
				},
			},
			{
				Name: "percona-xtradb",
				Versions: []string{
					"5.7",
				},
			},
		},
	}
}

func main() {
	catalog := CreateCatalogData()
	catalog.Sort()

	data, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
