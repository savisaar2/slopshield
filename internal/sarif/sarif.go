package sarif

import (
	"encoding/json"
	"io"
)

type Log struct {
	Version string `json:"version"`
	Schema  string `json:"$schema"`
	Runs    []Run  `json:"runs"`
}

type Run struct {
	Tool    Tool     `json:"tool"`
	Results []Result `json:"results"`
}

type Tool struct {
	Driver Driver `json:"driver"`
}

type Driver struct {
	Name           string `json:"name"`
	InformationURI string `json:"informationUri"`
	Rules          []Rule `json:"rules"`
}

type Rule struct {
	ID               string           `json:"id"`
	ShortDescription ShortDescription `json:"shortDescription"`
	FullDescription  ShortDescription `json:"fullDescription,omitempty"`
	HelpURI          string           `json:"helpUri,omitempty"`
	Help             Help             `json:"help,omitempty"`
}

type Help struct {
	Text     string `json:"text,omitempty"`
	Markdown string `json:"markdown,omitempty"`
}

type ShortDescription struct {
	Text string `json:"text"`
}

type Result struct {
	RuleID    string     `json:"ruleId"`
	Level     string     `json:"level,omitempty"`
	Message   Message    `json:"message"`
	Locations []Location `json:"locations"`
}

type Message struct {
	Text string `json:"text"`
}

type Location struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
}

type ArtifactLocation struct {
	URI string `json:"uri"`
}

func Generate(w io.Writer, hallucinations []string, sourceFile string) error {
	log := Log{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs: []Run{
			{
				Tool: Tool{
					Driver: Driver{
						Name:           "SlopShield",
						InformationURI: "https://github.com/savisaar2/slopshield",
						Rules: []Rule{
							{
								ID: "SLOP001",
								ShortDescription: ShortDescription{
									Text: "AI Hallucinated Package Detected",
								},
								FullDescription: ShortDescription{
									Text: "A dependency was found that does not exist in the official registry and is likely an AI hallucination.",
								},
								Help: Help{
									Text:     "AI hallucinated packages are a security risk. Attackers can register these names to perform supply chain attacks.",
									Markdown: "### Why this is a risk\nAI models frequently suggest non-existent packages. Attackers can register these names to perform supply chain attacks.\n\n### Recommendation\nVerify the package name and replace it with a reputable, existing alternative.",
								},
							},
						},
					},
				},
			},
		},
	}

	results := []Result{}
	for _, name := range hallucinations {
		results = append(results, Result{
			RuleID: "SLOP001",
			Level:  "error",
			Message: Message{
				Text: "Package '" + name + "' does not exist in the registry and might be an AI hallucination.",
			},
			Locations: []Location{
				{
					PhysicalLocation: PhysicalLocation{
						ArtifactLocation: ArtifactLocation{
							URI: sourceFile,
						},
					},
				},
			},
		})
	}
	log.Runs[0].Results = results

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log)
}
