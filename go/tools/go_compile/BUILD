package(default_visibility = ["//visibility:public"])

load("//go:def.bzl", "go_prefix", "go_library", "go_binary", "go_test")

SRCS = [
    "main.go",
    "compile.go",
]

# NOTE: we cannot use go_binary to bootstrap.
genrule(
    name = "bootstrap_go_compile",
    srcs = SRCS,
    outs = ["go_compile"],
    cmd = " && ".join([
        "export GOROOT=$$(dirname $(location //go/toolchain:go_tool))/..",
        ("$(location //go/toolchain:go_tool) tool compile -o $@.a -pack " +
         " -trimpath $$PWD $(SRCS)"),
        ("$(location //go/toolchain:go_tool) tool link -o $@ -buildmode exe" +
         " $@.a"),
        "rm -f $@.a",
    ]),
    executable = 1,
    message = "BootstrapGoCompile",
    output_to_bindir = 1,
    tools = [
        "//go/toolchain:go_tool",
    ],
)
