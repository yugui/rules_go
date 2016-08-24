workspace(name = "io_bazel_rules_go")

load("//go:def.bzl", "go_repositories", "go_internal_tools_deps")

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
    build_file_content = GLOG_BUILD,
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",
    remote = "https://github.com/golang/glog.git",
)

go_internal_tools_deps()

local_repository(
    name = "io_bazel_rules_go",
    path = ".",
)
