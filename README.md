# note-utils

A collection of utility functions meant for interacting with notes or other text-based content.
Mainly created this to integrate with ollama to use LLMs to summarize notes I take for work or other personal life.

## requirements

You need to have `ollama` installed locally, and ideally you have some of the popular models already installed so you don't have to wait a long time while they are pulled.
Models used in this repo:

- llama3.2:3b 
- deepseek-r1

## features

Currently supported features:

- summarizing text content (into markdown): `note-utils summarize` will take either a file (`--file`) or `stdin` and summarize the content for you, outputting it in markdown format.
- "cleaning up" text content (into markdown): `note-utils cleanup` will take `stdin` and output a "cleaned up" version, formatted in markdown. It aims to stay closer to the original source content, but improve readability.
