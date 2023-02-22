---
description: List Mongo Templates
---

# MongoTemplatesList

## Usage

```text
mongo templates list [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `templates` command:

```text
[t]
```

For `list` command:

```text
[l ls]
```

## Description

Retrieves a list of valid templates. These templates can be used to create MongoDB clusters; they contain properties, such as number of cores, RAM, and the storage size.

## Options

```text
      --cols strings   Set of columns to be printed on output 
                       Available columns: [TemplateId Name Edition Cores StorageSize Ram] (default [TemplateId,Name,Edition,Cores,StorageSize,Ram])
  -h, --help           help for list
      --no-headers     When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas mongo templates list
```

