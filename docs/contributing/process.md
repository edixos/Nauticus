## Contributing to Nauticus

Thank you for your interest in contributing to Nauticus. Whether it's a bug report, new feature, correction,
or additional documentation, we greatly value feedback and contributions from our community.

Please read through this document before submitting any issues or pull requests to ensure we have all the 
necessary information to effectively respond to your bug report or contribution.

## Submitting a Pull Request

This project uses the well-known pull request process from GitHub. To submit a
pull request, fork the repository and push any changes to a branch on the copy,
from there a pull request can be made in the main repo. Merging a pull request
requires the following steps to be completed before the pull request will
be merged:

* ideally, there is an issue that documents the problem or feature in depth.
* code must have a reasonable amount of test coverage (work in progress to integrate tests [#4](https://github.com/edixos/nauticus/issues/4))
* tests must pass (work in progress to integrate tests [#4](https://github.com/edixos/nauticus/issues/4))
* PR needs be reviewed and approved

Once these steps are completed the PR will be merged by a code owner.
We're using the pull request `assignee` feature to track who is responsible
for the lifecycle of the PR: review, merging, ping on inactivity, close.
We close pull requests or issues if there is no response from the author for
a period of time. Feel free to reopen if you want to get back on it.

## Proposal Process
Before we introduce significant changes to the project we want to gather feedback
from the community to ensure that we progress in the right direction before we
develop and release big changes. Significant changes include for example:

* creating new custom resources
* proposing breaking changes
* changing the behavior of the controller significantly

Please create a document in the `design/` directory based on the template `000-template.md`
and fill in your proposal. Open a pull request in draft mode and request feedback. Once the proposal is accepted and the pull request is merged we can create work packages and proceed with the implementation.
