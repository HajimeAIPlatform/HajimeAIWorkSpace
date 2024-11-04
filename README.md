# HajimeAIWorkSpace
The mono repo for HajimeAIWorkSpace

The workspace is managed by Bazel.

You can find bazel plugin on VSCode and Jetbrain IDEs.

We recommend you use VSCode and CLion as the developing tool.

Clion also support python.

You can try the demo c++ app / python app
```shell

bazel run //pythonp/apps/python_example_app:hello_hajime

```

Dependent Compilation

- python

pip-tools generates requirements.txt through requirements.in
```shell
bazel run //:requirements.update

```

The requirements_lock.txt will be updated. If you need to create a new version set of python
packages, you need to create a separate requirements.in file.

- golang

Add a New Dependency:
```shell
go get example.com/package

```
go get example.com/package
This will add example.com/package to your `go.mod` file.

Dependencies are defined in the `go.mod` file. To integrate these dependencies into the Bazel build system, execute the following command:
```shell
bazel run //:gazelle-update-repos

```
This command uses the Gazelle tool to read the go.mod file and update the deps.bzl file, which contains all the Go dependencies.

Coding style formatter for Python

We will be using yapf managing our code style.

https://github.com/google/yapf

for auto-formatting:
```shell
pip install yapf
yapf -r -i ./pythonp/apps/
```
-r represent recursive
-i means inplace

We will enforce yapf check in CI/CD.



