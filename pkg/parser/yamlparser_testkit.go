package parser

var instructionsYamlText = `
promptItems:
  - id: hello-world
    title: "Hello World"
    file: /path/to/hello_world.sh

  - id: goodbye-world
    title: "Goodbye World"
    file: /path/to/goodbye_world.sh

autoRun:
  - hello-world

autoCleanup:
  - goodbye-world
`
