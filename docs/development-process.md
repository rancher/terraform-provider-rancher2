# Rancher2 Provider

Terraform has recently evolved into a spotlight infrastructure provisioning tool that is maintained by Rancher engineering. As its popularity grows, so will its open source development.

Improvements to note since Jan 2023
* Rancher2 Provider docs
* Branching into `master 3.x.x` for Rancher 2.7 clusters and `release/v2 2.x.x` for Rancher 2.6 clusters
* Scripts to test Terraform RCs locally on unix and windows

### Development Process

If you are an open source contributor,

1. Add a feature / bug fix and smoke test it
2. Add or update tests
2. Open a PR with a description and how the feature / bug fix was tested
2. Get your PR merged
2. Add a test template with your test steps to the GitHub issue so the Terraform maintainers know how to test your solution.

```
# Test Template

Issue: <!-- link the issue or issues this PR resolves here -->
<!-- If your PR depends on changes from another pr link them here and describe why they are needed on your solution section. -->
  
# Problem
<!-- Describe the root cause of the issue you are resolving. This may include what behavior is observed and why it is not desirable. If this is a new feature describe why we need this feature and how it will be used. -->
  
# Solution
<!-- Describe what you changed to fix the issue. Relate your changes back to the original issue / feature and explain why this addresses the issue. -->
  
# Testing
<!-- Note: Confirm if the repro steps in the GitHub issue are valid, if not, please update the issue with accurate repro steps. -->
 
## Engineering Testing

### Manual Testing
<!-- Describe what manual testing you did (if no testing was done, explain why). -->
 
### Automated Testing
<!--If you added/updated unit/integration/validation tests, describe what cases they cover and do not cover. -->
 
## QA Testing Considerations
<!-- Highlight areas or (additional) cases that QA should test -->
  
### Regressions Considerations
<!-- Dedicated section to specifically call out any areas that with higher chance of regressions caused by this change, include estimation of probability of regressions -->
```
