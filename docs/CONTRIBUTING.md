# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue, email, or any other method with the owners of this repository before making a change.
Please note we have a [code of conduct](./CODE_OF_CONDUCT.md), please follow it in all your interactions with the project.

## Development environment setup

To set up a development environment, please follow these steps:

1. Clone the repo

   ```sh
   git clone https://github.com/Frolower/DiceDasher
   ```

2. Enter the folder

   ```sh
   cd DiceDasher
   ```

3. Run the docker containers

   ```sh
   make up
   ```
4. Stop the docker containers

   ```sh
   make down
   ```

## Issues and feature requests

You've found a bug in the source code, a mistake in the documentation or maybe you'd like a new feature? 
You can help us by [submitting an issue on GitHub](https://github.com/Frolower/DiceDasher/issues). 
Before you create an issue, make sure to search the issue archive -- your issue may have already been addressed!

Please try to create bug reports that are:

- _Reproducible._ Include steps to reproduce the problem.
- _Specific._ Include as much detail as possible: which version, what environment, etc.
- _Unique._ Do not duplicate existing opened issues.
- _Scoped to a Single Bug._ One bug per report.

**Even better: Submit a pull request with a fix or new feature!**

### How to submit a Pull Request

1. Search our repository for open or closed
   [Pull Requests](https://github.com/Frolower/DiceDasher/pulls)
   that relate to your submission. You don't want to duplicate effort.
2. Fork the project
3. Create your feature branch (`git checkout -b feature/amazing_feature`) DiceDasher is using type-prefixed branch naming [convention](https://conventional-branch.github.io/).
4. Commit your changes (`git commit -m 'add amazing_feature'`) DiceDasher doesn't use any strict conventions for the commit messages, it is up to you to keep them reasonable. Commits are suggested to be imperative, short, present tense
5. Push to the branch (`git push origin feature/amazing_feature`)
6. [Open a Pull Request](https://github.com/Frolower/DiceDasher/compare?expand=1)
