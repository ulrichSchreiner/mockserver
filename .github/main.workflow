workflow "New Docker image" {
  on = "push"
  resolves = ["Publish Image"]
}

action "Docker Login" {
    uses = "actions/docker/login@master"
    secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Publish Image" {
    needs = ["Docker Login"]
    uses = "./.github/devaction/"
    args = ["make", "build", "push"]
}
