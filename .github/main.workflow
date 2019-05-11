workflow "New Docker image" {
  on = "push"
  resolves = ["Publish Image"]
}

action "Docker Login" {
    uses = "actions/docker/login@master"
    secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Build Image" {
    needs = ["Docker Login"]
    uses = "./.github/devaction/"
    args = ["make", "build"]
}

action "Publish Image" {
    needs = ["Build Image"]
    uses = "./.github/devaction/"
    args = ["make", "push"]
}
