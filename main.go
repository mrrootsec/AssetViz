package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	tld "github.com/jpillora/go-tld"
)

// DomainTree represents a tree structure for domains
type DomainTree map[string]DomainTree

// main function handles command-line arguments and file processing
func main() {
	var filePath string
	flag.StringVar(&filePath, "f", "", "Path to the file containing subdomain names")
	flag.Parse()

	if filePath == "" {
		// Process input from stdin if no file path provided
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			processInput(os.Stdin)
		} else {
			fmt.Println("Usage: assetviz -f filename OR provide input via stdin")
			return
		}
	} else {
		// Process input from file
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()
		processInput(file)
	}
}

// isValidDomain checks if a domain is valid
func isValidDomain(domain string) bool {
	u, err := tld.Parse("http://" + domain)
	if err != nil {
		return false
	}
	return u.Domain != ""
}

// processInput reads input from a file or stdin and builds the domain tree
func processInput(input *os.File) {
	domainTree := make(DomainTree)
	scanner := bufio.NewScanner(input)
	var encounteredError bool

	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())

		// Check for empty lines and single-dot domains
		if domain != "" && domain != "." {
			// Trim protocol, trailings and port from the domain name
			domain = strings.TrimPrefix(domain, "http://")
			domain = strings.TrimPrefix(domain, "https://")
			domain = strings.ReplaceAll(domain, "..", ".")
			domain = strings.TrimSuffix(domain, ".")
			domain = strings.Split(domain, ":")[0]

			// Validate and update domain tree
			if isValidDomain(domain) {
				updateDomainTree(domainTree, domain)
			} else {
				if !encounteredError {
					fmt.Println("File contains invalid input")
					encounteredError = true
					return
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		return
	}

	// Convert domain tree to JSON
	jsonBytes, err := json.MarshalIndent(domainTree, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}
	generateHTMLReport(jsonBytes)
}

// updateDomainTree updates the domain tree with the given domain
func updateDomainTree(domainTree DomainTree, domain string) {
	parts := strings.Split(domain, ".")
	tldIndex := len(parts) - 1
	currentTree := domainTree

	for i := len(parts) - 1; i >= 0; i-- {
		subdomain := parts[i]
		fullDomain := strings.Join(parts[i:], ".")

		if i == tldIndex {
			if _, exists := currentTree[subdomain]; !exists {
				currentTree[subdomain] = make(DomainTree)
			}
			currentTree = currentTree[subdomain]
		} else {
			if _, exists := currentTree[fullDomain]; !exists {
				currentTree[fullDomain] = make(DomainTree)
			}
			currentTree = currentTree[fullDomain]
		}
	}
}

func generateHTMLReport(jsonData []byte) {
	// Create a timestamp for the report
	currentTime := time.Now().Format("2006-01-02_15-04-05")

	// Define the file path for the report
	reportFileName := fmt.Sprintf(".report/assetviz_report_%s.html", currentTime)

	// Create the .report directory if it doesn't exist
	err := os.MkdirAll(".report", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating .report directory:", err)
		return
	}

	// Define the HTML template
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>AssetViz</title>
    <link type="text/css" rel="stylesheet" href="https://cdn.jsdelivr.net/npm/jsmind@0.8.1/style/jsmind.css" />
    <link rel="icon" type="image/x-icon" href="/img/favicon.ico">
	<script src="https://cdn.jsdelivr.net/npm/jsmind@0.8.1/es6/jsmind.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/jsmind@0.8.1/es6/jsmind.draggable-node.js"></script>
	<script type=text/javascript>
        var jsondata = {{.}}

        function convertJSONToJSMindData(myjson) {
            let mindMapTree = {
                meta: {
                    name: 'AssetViz',
                    author: 'mrrootsec',
                    version: '1.0',
                },
                format: 'node_tree',
                data: {
                    id: 'root',
                    topic: 'AssetViz',
                    children: []
                }
            }

            Object.entries(myjson).forEach(([key, value]) => {
                mindMapTree.data.children.push(CreateJSMindNode(key, value))
            });

            console.log(mindMapTree)
            return mindMapTree
        }

        function CreateJSMindNode(title, children) {
            let node = {
                id: title,
                topic: title,
                children: []
            }
            if (Object.keys(children).length > 0) {
                Object.entries(children).forEach(([key, value]) => {
                    node.children.push(CreateJSMindNode(key, value))
                });
            }
            return node
        }
    </script>
    <style type="text/css" rel="stylesheet">
        body {
            margin: 0;
            padding: 0;
            width: 100%;
            height: 100vh;
            background: rgb(32, 32, 32);
            background-image: radial-gradient(rgba(255, 255, 255, 0.1) 1px, transparent 1px);
            background-size: 20px 20px; /* Adjust the size of the grid dots */
        }

        footer {
            color: white;
            text-align: center;
            position: inherit;
            bottom: 0;
            width: 100%;
        }

        #jsmind_container {
            position: left;
            padding: 0;
            margin: 0;
            width: 100%;
            height: 100%;
        }

        .jsmind-inner {
            position: relative;
            overflow: auto;
            width: 100%;
            height: 100%;
            outline: none;
        }

        .jsmind-inner canvas {
            position: absolute;
        }

        /* default theme */
        jmnode {
            padding: 10px;
            background-color: #2b2c3e;
            color: #7fc0f1;
            border-radius: 5px;
            box-shadow: 1px 1px 1px #666;
            font: 16px/1.125 Verdana, Arial, Helvetica, sans-serif;
        }

        jmnode:hover {
            box-shadow: 2px 2px 8px #7fc0f1;
            color: #333;
        }

        jmnode.selected {
            background-color: rgb(111, 111, 235);
            color: #fff;
            box-shadow: 2px 2px 8px #5478da;
        }

        jmnode.root {
            font-size: 24px;
            font-weight: bold;
        }

        jmexpander {
            border-color: #7fc0f1;
        }

        jmexpander:hover {
            border-color: #255ec7;
            background-color: yellow;
            transition: background-color ease-out 100ms
        }

        @media screen and (max-device-width: 1024px) {
            jmnode {
                padding: 5px;
                border-radius: 3px;
                font-size: 14px;
            }

            jmnode.root {
                font-size: 21px;
            }
        }
    </style>
</head>
<body>
    <div id="jsmind_container" class="container-fluid"></div>
    <footer>
        <p style="font-weight:bold">AssetViz V1.0 &middot; made with <span style="color:red;font-weight:bold"> ❤️ </span> by <a style="color:white" href="https://twitter.com/_mohd_saqlain" target="_blank">Mrootsec</a></p>
    </footer>
    <script type="text/javascript">
        function load_jsmind() {
            var options = {
                container: 'jsmind_container',
                editable: false,
                view: {
                    engine: 'svg', // engine for drawing lines between nodes in the mindmap
                    hmargin: 200, // Minimum horizontal distance of the mindmap from the outer frame of the container
                    vmargin: 200, // Minimum vertical distance of the mindmap from the outer frame of the container
                    line_width: 1, // thickness of the mindmap line
                    line_color: '#7fc0f1', // Thought mindmap line color
                    line_style: 'curved', // line style, straight or curved
                    draggable: true, // Drag the mind map with your mouse, when it's larger that the container
                    hide_scrollbars_when_draggable: false, // Hide container scrollbars, when mind map is larger than container and draggable option is true.
                    node_overflow: 'hidden' // Text overflow style in node 
                },
                layout: {
                    hspace: 50, // Horizontal spacing between nodes
                    vspace: 20, // Vertical spacing between nodes
                    pspace: 13, // Horizontal spacing between node and connection line (to place node expander)
                    cousin_space: 5 // Additional vertical spacing between child nodes of neighbor nodes
                },
                shortcut: {
                    enable: true, // whether to enable shortcut
                    handles: {}, // Named shortcut key event processor
                    mapping: { // shortcut key mapping
                        toggle: 32, // <Space>
                        left: 37, // <Left>
                        up: 38, // <Up>
                        right: 39, // <Right>
                        down: 40, // <Down>
                    }
                }
            };
            var jm = new jsMind(options);
            jm.show(convertJSONToJSMindData(jsondata));
        }
        load_jsmind();
    </script>
</body>
</html>
`
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		fmt.Println("Error parsing HTML template:", err)
		return
	}

	// Create or open the report file
	reportFile, err := os.Create(reportFileName)
	if err != nil {
		fmt.Println("Error creating report file:", err)
		return
	}
	defer reportFile.Close()

	// Execute the template and write to the report file
	err = tmpl.Execute(reportFile, string(jsonData))
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Printf("HTML report generated: %s\n", reportFileName)
}
