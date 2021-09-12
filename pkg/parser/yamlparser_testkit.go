package parser

var anchorFolderInfoYamlText = `
type: application
name: app
command:
  use: app
  short: "Application commands"
`

var instructionsOnlyActionsYamlText = `
instructions:
  actions:
    - id: hello-world
      title: "Hello World"
      description: "Print 'Hello World' to stdout"
      script: |
        echo "Hello"
        echo "World"

    - id: goodbye-world
      title: "Goodbye World"
      scriptFile: app/first-app/sub-folder/goodbye_world.sh

    - id: global-hello-world
      title: "Global Hello World"
      scriptFile: k8s/scripts/global_hello_world.sh
`

var instructionsWithWorkflowsYamlText = `
instructions:
  actions:
    - id: hello-world
      title: "Hello World"
      description: "Print 'Hello World' to stdout"
      script: |
        echo "Hello"
        echo "World"

    - id: goodbye-world
      title: "Goodbye World"
      scriptFile: app/first-app/sub-folder/goodbye_world.sh

    - id: global-hello-world
      title: "Global Hello World"
      scriptFile: k8s/scripts/global_hello_world.sh

  workflows:
    - id: talk-to-the-world
      description: "Talk to the world, starting from hello up to goodbye"
      tolerateFailures: false
      actionIds:
        - hello-world
        - goodbye-world
`
