# slack-cleaner
Tool to clean slack messages

# Installation
1. Download version [slack-cleaner](https://github.com/ypelud/slack-cleaner)
2. Copy config.toml.sample to config.toml
3. Replace slackId with your token see [Legacy token](https://get.slack.help/hc/en-us/articles/215770388)

# Usage
All examples are on channel general

Delete all messages
```shell
./slack-cleaner general
```

Delete all messages before 01/01/2018
```shell
./slack-cleaner -date=20180101 general
```

Delete all messages from user bob
```shell
./slack-cleaner -user=bob general
```

Delete all messages from user bob with dry-run
```shell
./slack-cleaner -user=bob -n general
```


