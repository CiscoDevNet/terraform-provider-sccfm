[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8023/badge)](https://www.bestpractices.dev/projects/8023)

# SCC Firewall Manager Terraform Provider

This repo provides a terraform provider to provision resources on the Cisco SCC Firewall Manager platform.

## Structure

We make use of Go workspaces - this is to split the responsibility of the provider and the Go SCC Firewall Manager client. 
Eventually the client will be moved to its own repo, but in the interest of hitting the ground running they both live in here.

```
.
├── client        # Golang SCC Firewall Manager client - currently using the module name github.com/CiscoDevnet/terraform-provider-sccfm/go-client
├── provider      # Terraform provider
└── README.md
```

## Requirements

* Go (1.20)
  - macos install: `brew install go`
* tfenv 
  - macos install: `brew install tfenv`

### Gitleaks

We use [Gitleaks](https://github.com/gitleaks/gitleaks) to catch secrets being committed to the repository by accident. The first line of defense is before you ever push to GitHub using a pre-commit hook.

Please enable the pre-commit hook before you commit anything to this repository, even in a branch.

- Install `pre-commit`
```
brew install pre-commit
```
- Update the pre-commit configuration:
```
pre-commit autoupdate
```
- Install the pre-commit configuration as a pre-commit hook to your local Git repo:
```
pre-commit install
```

Now any commits you make will be scanned by Gitleaks


## Acceptance Tests

**Acceptance tests will create real resources!**

Ensure you have met the requirements.

Ensure you have an active terraform:

```bash
tfenv use 1.3.1
```

Then in the provider dir you can run the acceptance tests:

```bash
cd ./provider
ACC_TEST_CISCO_CDO_API_TOKEN=<CDO_API_TOKEN> make testacc
```

## Linting

Ensure you have golangci-lint installed:
```
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0
```

Run following command in the `client` or `provider` directory.
```bash
~/go/bin/golangci-lint run #change ~/go/bin to your $GOBIN
```

## Running Examples
Examples are provided so that you can do the usual `plan`, `apply`, `destroy` etc in the folders under `provider/examples` directory.
#### Setup
1. Build provider locally.
   1. Figure out where your executable will be built.
      ```bash
      go env GOBIN
      ```
      **Note:** If empty output, it is default to `~/go/bin`
   2. Build executable.
      ```bash
      cd provider
      go install .
      ```
   3. Verify installation
      The installation path that you figured out above should contain an executable called **terraform-provider-sccfm**.
2. Tell terraform to use your local build
   1. Modify local terraform configuration
      ```bash
      vim ~/.terraformrc
      ```
   2. Copy the following into it:
      ```terraform
      provider_installation {
   
        dev_overrides {
            "hashicorp.com/CiscoDevnet/sccfm" = "<DIRECTORY_OF_YOUR_EXECUTABLE>"
        }
  
        # For all other providers, install them directly from their origin provider
        # registries as normal. If you omit this, Terraform will _only_ use
        # the dev_overrides block, and so no other providers will be available.
        direct {}
      }
      ```
## Running
1. Navigate to a folder under `provider/examples`, e.g. `provider/examples/resources/asa`.
2. Run
   ```bash
   terraform plan
   ```
3. Output (sample)
   ```
    |
    │ Warning: Provider development overrides are in effect
    │
    │ The following provider development overrides are set in the CLI configuration:
    │  - hashicorp.com/CiscoDevnet/sccfm in /Users/weilluo/go/bin
    │
    │ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
    |
    ... rest
   ```

## Regenerating docs

If you make any changes to the resources and data sources provided by this provider, you will need to regenerate the docs, otherwise the Github actions triggered by this pull request will fail. To do this, run:
```
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name sccfm --rendered-provider-name "SCC Firewall Manager Provider" --rendered-website-dir ../docs
```

## Releasing Manually
> **Note**: Tags and releases are automatically done in GitHub action when pushed/merged PR into main branch and checks passed, there is no need for releasing manually in normal scenario.

To release a new version of the Terraform SCC Firewall Manager Provider manually, perform the following steps.

- Checkout main branch.
- List current available tags: `git tag`.
- Add a tag: `git tag vMAJOR.MINOR.PATCH`, e.g. `git tag v0.1.3` (following semver conventions as described in https://www.semver.org).
  - To add a tag for past commit: `git tag -a vMAJOR.MINOR.PATCH COMMIT_HASH`, e.g. `git tag -a v1.2.3 9fceb02`. 
- Push the tag: `git push --tags origin main`.

## Troubleshooting
- Error: Inconsistent dependency lock file
  ```
  provider hashicorp/CiscoDevnet/sccfm: required by this configuration but no version is selected
  ```
  - This means you have not setup the dev override properly, make sure your `~/.terraformrc` has the right override for the provider in question.

## Gitleaks License

The Gitleaks License is free, and stored in the `GITLEAKS_LICENSE` secret. In addition, it is saved to [Conjur](https://secrets.cisco.com/conjur/nonprod/eng/cdo/gitleaks-license). Speak to Jay, Doron, Siddhu, or Pedro to access it.

## Running acceptance test locally against Github environment
1. Create a file called `./provider/.github-action.local.env` locally
2. File should contain the following:
   ```bash
   ASA_RESOURCE_SDC_PASSWORD=<fill this>
   IOS_RESOURCE_PASSWORD=<fill this>
   DUO_ADMIN_PANEL_RESOURCE_INTEGRATION_KEY=<fill this>
   DUO_ADMIN_PANEL_RESOURCE_SECRET_KEY=<fill this>
   ```
3. Login to the github environment in ci, the tenant name is `terraform-provider-cdo@lockhart.io`
4. Generate a api token
5. In the console, export it as environment variable, named `ACC_TEST_CISCO_CDO_API_TOKEN`
6. Run `cd ./provider`
7. Run `make testacc`
