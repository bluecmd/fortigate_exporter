# Contributing

Fortigate exporter uses GitHub to manage reviews of pull requests.

* If you are a new contributor see: [Steps to Contribute](#steps-to-contribute)

* If you have a trivial fix or improvement, go ahead and create a pull request,
  addressing it to @bluecmd or @secustor.

* If you plan to do something more impacting, first discuss your ideas
  in our [chat room](https://matrix.to/#/#fortigate_exporter:matrix.org) or the [discussions](https://github.com/bluecmd/fortigate_exporter/discussions/landing) page.
  This will avoid unnecessary work and surely give you and us a good deal
  of inspiration.

* Relevant coding style guidelines are the [Go Code Review
  Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments)

* Make sure you understand how Prometheus works [Prometheus Getting started](https://prometheus.io/docs/prometheus/latest/getting_started/)

* Make sure created metrics do respect [Prometheus naming guidelines](https://prometheus.io/docs/practices/naming/)

* Make sure your code respects the [Prometheus implementation guidelines](https://prometheus.io/docs/practices/instrumentation/#things-to-watch-out-for)

## Steps to Contribute

Should you wish to work on an issue, please claim it first by commenting on the GitHub issue that you want to work on it. This is to prevent duplicated efforts from contributors on the same issue.

To setup your dev environment you need:
- A way to run Makefiles, either Linux or WSL is recommended
- latest golang binaries (see [Golang Install](https://golang.org/doc/install))

For compiling and testing your changes do:
```
# For building.
make

# For checking code.
make vet
make fmt-fix

# For testing.
make test         # Make sure all the tests pass before you commit and push :)
```

## Pull Request Checklist

* Branch from the main branch and, if needed, rebase to the current main branch before submitting your pull request. If it doesn't merge cleanly with main you may be asked to rebase your changes.

* Commits should be as small as possible, while ensuring that each commit is correct independently (i.e., each commit should compile and pass tests).

* Commits messages have to adhere to the [conventional commits specification](https://www.conventionalcommits.org/en/v1.0.0/#summary).

* If multiple commits address independent things, use multiple PRs - avoid creating one PR that fix two or more things. 

* If your patch is not getting reviewed or you need a specific person to review it, you can @-reply a reviewer asking for a review in the pull request or a comment, or you can ask for a review in our [chat room](https://matrix.to/#/#fortigate_exporter:matrix.org).

* Add tests relevant to the fixed bug or new feature.

## Naming convention

New metrics name should follow the [Prometheus naming guidelines](https://prometheus.io/docs/practices/naming/) and [Prometheus implementation guidelines](https://prometheus.io/docs/practices/instrumentation/#things-to-watch-out-for)

**NEVER** change the metric format in a backwards incompatible way (e.g. rename metrics, remove labels) if the metric has been merged to main. Once a metric has been merged to main, it will stay there unless some extraordinary happens and we need to remove it. Our users rely on that the metrics remain stable, so we strive to keep it that way. If we need to drop metrics or do larger refactors those should be queued to the next major version bump (i.e. v2.0.0 or equivalent).

Category name used to include/exclude probes in the exporter configuration file should, as much as possible, follow those rules:
- By default, probe category name should be derived from the API URL path following the `api/v2/monitor` part (ex: `api/v2/monitor/user/fsso` become `User/Fsso`)
- A probe category name should never be a subpart of another probe category name, in order to keep inclusion/exclusion declaration simple and usable.
- probe category name should always match : `^([a-ZA-Z0-9]+/)+[a-ZA-Z0-9]+$`
- Removing a prefix part is possible if not conflicting or bringing ambiguity with current and futur probe category name/URL path (ex: `BGP/Neighbo...` instead of `Router/BGP/Neighbo...` -- shorter is better
- Adding a suffix part to avoid prefix matching unwilled inclusion/exclusion is possible (ex: `api/v2/monitor/vpn/ssl` and `api/v2/monitor/vpn/ssl/stats` goes to `VPN/Ssl/Connections` and `VPN/Ssl/Stats`)
- Creating an additional category section not present in the API URL path is possible to make grouping more efficient and anticipate present or futur Fortigate API prefix match issues (ex: `monitor/system/time` conflict with `monitor/system/timezone`, then category name can be `System/Time/Clock` and `System/Time/Timezone`)
