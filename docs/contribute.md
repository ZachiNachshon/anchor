<h1 id="contribute" align="center">Contribute<br><br></h1>

- [Guidelines](#guidelines)
- [Testing Locally](#testing-locally)

<br>

<h3 id="guidelines">Guidelines</h3>

- PRs need to have a clear description of the problem they are solving
- PRs should be small
- Code without tests is not accepted, PRs must not reduce tests coverage
- Contributions must not add additional dependencies
- Before creating a PR, make sure your code is well formatted, abstractions are named properly and design is simple
- In case your contribution can't comply with any of the above please start a GitHub issue for discussion

<br>

<h3 id="testing-locally">Testing Locally</h3>

This repository had been developed using the TDD methodology (Test Driven Development).

Tests allow you to make sure your changes work as expected, don't break existing code and keeping test coverage high.

Running tests locally allows you to have short validation cycles instead of waiting for the PR status to complete.

**How to run a test suite?**

1. Clone the `anchor` repository
2. Run `make run-tests` to use the locally installed Go runtime
3. Alternatively, run `make run-tests-containerized` to use the same Go runtime which is supported by this repository

   | :bulb: Note |
   | :--------------------------------------- |
   | `tparse` should be installed for performing `go test` output analysis ([instructions in here](https://github.com/mfridman/tparse))|
