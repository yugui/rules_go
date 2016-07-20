load("//go:def.bzl", "go_binary")

go_binary(
    name = "fetch_repo",
    srcs = [
        "main.go",
        "repo.go",
    ],
    deps = [
        "@com_github_masterminds_vcs//:go_default_library",
        "@org_golang_x_tools//:go/vcs",
    ],
)
