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

```shell
bazel run //:requirements.update

```

The requirements_lock.txt will be updated. If you need to create a new version set of python
packages, you need to create a separate requirements.in file.



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



