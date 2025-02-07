# Agent Town

## Bazel 8

This program is designed to only work with best practices for Bazel 8.

Bazel 8 is not compatible with Bazel 7.xor lower.

We enforce using bzlmod and rely on the dependencies to be managed by MODULES.bzl in root file.

Currently most bazel plugins are not compatible with Bazel 8. But eventually they will be.

## Introduction

Agent Town is designed to simulate a group of agent that can interact with each other using message chan and task chan.
The system is designed to be general purpose and can be used to simulate any kind of actions. Agents can be extended to
be powered by AI or any other kind of logic.

A agent can receive two types of info:

1. Message: A message is a simple string that can be sent to an agent. The agent can choose to ignore it or act on it.
2. Task: A task is a more complex object that can be sent to an agent. The agent can choose to ignore it or act on it.
3. Agent can also send message and task to other agents.
4. Task can also trigger other tasks and messages.

Currently due to the transition to Bazel 8, the logging system is not yet merged. Eventually we will merge all code to be
using Bazel 8 for code management.

## FAQ

1. Can I use older version of Bazel?
    No
