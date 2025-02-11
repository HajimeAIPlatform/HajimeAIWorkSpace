# HajimeAIWorkSpace
The mono repo for HajimeAIWorkSpace

The workspace is managed by Bazel.

You can find bazel plugin on VSCode and Jetbrain IDEs.

We recommend you use VSCode and CLion as the developing tool.

Clion also support python.

## Recommand GO practice
### Run application
```shell
bazel run //golangp/apps/hajime_center
```

### Add external dependencies
Please regenerate BUILD file to ensure compile success.
```shell
bazel run @rules_go//go get example.com/package
```

### Auto (re)generate BUILD file
In most cases, manual editing of the BUILD file is not necessary.  
It is recommanded to use the command to automatically update dependencies.  
```shell
bazel run //:gazelle
```

### Run testcase
```shell
bazel test --test_filter=TestSmokeTest //golangp/apps/AgentTown/test:common_test
```

### Recommand VScode plugin
- https://marketplace.visualstudio.com/items?itemName=NVIDIA.bluebazel

### Reference
https://github.com/bazel-contrib/rules_go/blob/master/docs/go/core/bzlmod.md

## Python
### Dependent Compilation

- python

pip-tools generates requirements.txt through requirements.in
```shell
bazel run //:requirements.update

```

The requirements_lock.txt will be updated. If you need to create a new version set of python
packages, you need to create a separate requirements.in file.

### Coding style formatter for Python

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