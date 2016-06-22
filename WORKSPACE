workspace(name = "io_bazel_rules_go")

load("//go:def.bzl", "go_repositories")
go_repositories()

GLOG_BUILD = """
load("@//go:def.bzl", "go_prefix", "go_library")

go_prefix("github.com/golang/glog")

go_library(
    name = "go_default_library",
    srcs = [
        "glog.go",
        "glog_file.go",
    ],
    visibility = ["//visibility:public"],
)
"""

new_git_repository(
    name = "com_github_golang_glog",
    remote = "https://github.com/golang/glog.git",
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",
    build_file_content = GLOG_BUILD,
)

git_repository(
    name = "io_bazel_buildifier",
    remote = "https://github.com/bazelbuild/buildifier.git",
    commit = "84cdc95dd453430af1206c1bfc9e4cddb45e7670",
)

# TODO(yugui) DO NOT SUBMIT. Remove this reference before sending PR.
local_repository(
    name = "io_bazel_rules_go",
    path = "/Users/yugui/dev/bazel/rules_go",
)

X_TOOLS_BUILD = """
load("@//go:def.bzl", "go_prefix", "go_library")

go_prefix("golang.org/x/tools")

go_library(
    name = "go/vcs",
    srcs = glob(
      include = ["go/vcs/*.go"],
      exclude = ["go/vcs/*_test.go"],
    ),
    visibility = ["//visibility:public"],
)
"""

new_git_repository(
    name = "org_golang_x_tools",
    remote = "https://github.com/golang/tools.git",
    commit = "a2a552218a0e22e6fb22469f49ef371b492f6178",
    build_file_content = X_TOOLS_BUILD,
)
