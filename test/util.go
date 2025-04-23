package tests

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	ec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/go-github/v53/github"
	aws "github.com/gruntwork-io/terratest/modules/aws"
	g "github.com/gruntwork-io/terratest/modules/git"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"golang.org/x/oauth2"
)

func GetRancherReleases() (string, string, string, error) {
	releases, err := getRancherReleases()
	if err != nil {
		return "", "", "", err
	}
	versions := filterPrerelease(releases)
	if len(versions) == 0 {
		return "", "", "", errors.New("no eligible versions found")
	}
	zeroPadVersionNumbers(&versions)
	sortVersions(&versions)
	filterDuplicateMinors(&versions)
	removeZeroPadding(&versions)
	latest := versions[0]
	stable := latest
	lts := stable
	if len(versions) > 1 {
		stable = versions[1]
	}
	if len(versions) > 2 {
		lts = versions[2]
	}
	return latest, stable, lts, nil
}

func getRancherReleases() ([]*github.RepositoryRelease, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		fmt.Println("GITHUB_TOKEN environment variable not set")
		return nil, errors.New("GITHUB_TOKEN environment variable not set")
	}

	// Create a new OAuth2 token using the GitHub token
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tokenClient := oauth2.NewClient(context.Background(), tokenSource)

	// Create a new GitHub client using the authenticated HTTP client
	client := github.NewClient(tokenClient)

	var releases []*github.RepositoryRelease
	releases, _, err := client.Repositories.ListReleases(context.Background(), "rancher", "rancher", &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func GetRke2Releases() (string, string, string, error) {
	releases, err := getRke2Releases()
	if err != nil {
		return "", "", "", err
	}
	versions := filterPrerelease(releases)
	if len(versions) == 0 {
		return "", "", "", errors.New("no eligible versions found")
	}
	zeroPadVersionNumbers(&versions)
	sortVersions(&versions)
	filterDuplicateMinors(&versions)
	removeZeroPadding(&versions)
	latest := versions[0]
	stable := latest
	lts := stable
	if len(versions) > 1 {
		stable = versions[1]
	}
	if len(versions) > 2 {
		lts = versions[2]
	}
	return latest, stable, lts, nil
}

func getRke2Releases() ([]*github.RepositoryRelease, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		fmt.Println("GITHUB_TOKEN environment variable not set")
		return nil, errors.New("GITHUB_TOKEN environment variable not set")
	}

	// Create a new OAuth2 token using the GitHub token
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tokenClient := oauth2.NewClient(context.Background(), tokenSource)

	// Create a new GitHub client using the authenticated HTTP client
	client := github.NewClient(tokenClient)

	var releases []*github.RepositoryRelease
	releases, _, err := client.Repositories.ListReleases(context.Background(), "rancher", "rke2", &github.ListOptions{})
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func filterReleaseCandidate(v *[]string) {
	var fv []string
	versions := *v
	for i := 1; i < len(versions); i++ {
		if strings.Contains(versions[i], "-") != true {
			fv = append(fv, versions[i])
		}
	}
	*v = fv
}

func filterPrerelease(r []*github.RepositoryRelease) []string {
	var versions []string
	for _, release := range r {
		version := release.GetTagName()
		if !release.GetPrerelease() {
			versions = append(versions, version)
			// [
			//   "v1.28.14+rke2r1",
			//   "v1.30.1+rke2r3",
			//   "v1.29.4+rke2r1",
			//   "v1.30.1+rke2r2",
			//   "v1.29.5+rke2r2",
			//   "v1.30.1+rke2r1",
			//   "v1.27.20+rke2r1",
			//   "v1.30.0+rke2r1",
			//   "v1.29.5+rke2r1",
			//   "v1.28.17+rke2r1",
			// ]
		}
	}
	return versions
}

func sortVersions(v *[]string) {
	slices.SortFunc(*v, func(a, b string) int {
		return cmp.Compare(b, a)
		//[
		//  v1.4.1+rke2r3,
		//  v1.30.1+rke2r3,
		//  v1.30.1+rke2r2,
		//  v1.30.1+rke2r1,
		//  v1.30.0+rke2r1,
		//  v1.29.5+rke2r2,
		//  v1.29.5+rke2r1,
		//  v1.29.4+rke2r1,
		//  v1.28.17+rke2r1,
		//  v1.28.14+rke2r1,
		//  v1.27.20+rke2r1,
		//]
	})
}

func filterDuplicateMinors(v *[]string) { // assumes versions are sorted already
	var fv []string
	versions := *v
	fv = append(fv, versions[0])
	for i := 1; i < len(versions); i++ {
		c := versions[i]
		p := versions[i-1]
		cp := strings.Split(c[1:], "+") //["1.30.1","rke2r3"]
		pp := strings.Split(p[1:], "+") //["1.30.1","rke2r2"]
		if cp[0] != pp[0] {
			cpp := strings.Split(cp[0], ".") //["1","30","1]
			ppp := strings.Split(pp[0], ".") //["1","30","1]
			if cpp[1] != ppp[1] {
				fv = append(fv, c)
				//[
				//  v1.30.1+rke2r3,
				//  v1.29.5+rke2r2,
				//  v1.28.17+rke2r1,
				//  v1.27.20+rke2r1,
				//]
			}
		}
	}
	*v = fv
}

func zeroPadVersionNumbers(v *[]string) {
	var zv []string
	versions := *v
	for i := 0; i < len(versions); i++ {
		vp := strings.Split(versions[i], "+") //["v1.3.1","rke2r3"] OR ["v2.5.4"] if no "+"
		vpp := strings.Split(vp[0], ".")      //["v1","3","1]
		major := vpp[0]                       // assumes single digit major
		minor := ""
		trivial := ""
		if len(vpp[1]) < 2 {
			minor = fmt.Sprintf("0%s", vpp[1]) // assumes double digit versions
		} else {
			minor = vpp[1]
		}
		if len(vpp[2]) < 2 {
			trivial = fmt.Sprintf("0%s", vpp[2]) // assumes double digit versions
		} else {
			trivial = vpp[2]
		}
		if len(vp) > 1 {
			version := fmt.Sprintf("%s.%s.%s+%s", major, minor, trivial, vp[1]) //"v1.03.01+rke2r3"
			zv = append(zv, version)
		} else {
			version := fmt.Sprintf("%s.%s.%s", major, minor, trivial) //"v1.03.01"
			zv = append(zv, version)
		}
	}
	*v = zv
}

func removeZeroPadding(v *[]string) {
	var zv []string
	versions := *v
	for i := 0; i < len(versions); i++ {
		vp := strings.Split(versions[i], "+") //["v1.03.01","rke2r3"] OR ["v2.05.04"] if no "+"
		vpp := strings.Split(vp[0], ".")      //["v1","03","01]
		major := vpp[0]                       // assumes single digit major
		minor := vpp[1]
		trivial := vpp[2]
		if minor[0] == '0' {
			minor = minor[1:]
		}
		if trivial[0] == '0' {
			trivial = trivial[1:]
		}
		if len(vp) > 1 {
			version := fmt.Sprintf("%s.%s.%s+%s", major, minor, trivial, vp[1]) //"v1.3.1+rke2r3"
			zv = append(zv, version)
		} else {
			version := fmt.Sprintf("%s.%s.%s", major, minor, trivial) //"v1.3.1"
			zv = append(zv, version)
		}
	}
	*v = zv
}

func CreateKeypair(t *testing.T, region string, owner string, id string) (*aws.Ec2Keypair, error) {
	t.Log("Creating keypair...")
	// Create an EC2 KeyPair that we can use for SSH access
	keyPairName := fmt.Sprintf("terraform-ci-%s", id)
	keyPair := aws.CreateAndImportEC2KeyPair(t, region, keyPairName)

	// tag the key pair so we can find in the access module
	client, err := aws.NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	k := "key-name"
	keyNameFilter := ec2.Filter{
		Name:   &k,
		Values: []*string{&keyPairName},
	}
	input := &ec2.DescribeKeyPairsInput{
		Filters: []*ec2.Filter{&keyNameFilter},
	}
	result, err := client.DescribeKeyPairs(input)
	if err != nil {
		return nil, err
	}

	err = aws.AddTagsToResourceE(t, region, *result.KeyPairs[0].KeyPairId, map[string]string{"Name": keyPairName, "Owner": owner})
	if err != nil {
		return nil, err
	}

	// Verify that the name and owner tags were placed properly
	k = "tag:Name"
	keyNameFilter = ec2.Filter{
		Name:   &k,
		Values: []*string{&keyPairName},
	}
	input = &ec2.DescribeKeyPairsInput{
		Filters: []*ec2.Filter{&keyNameFilter},
	}
	result, err = client.DescribeKeyPairs(input)
	if err != nil {
		return nil, err
	}

	k = "tag:Owner"
	keyNameFilter = ec2.Filter{
		Name:   &k,
		Values: []*string{&owner},
	}
	input = &ec2.DescribeKeyPairsInput{
		Filters: []*ec2.Filter{&keyNameFilter},
	}
	result, err = client.DescribeKeyPairs(input)
	if err != nil {
		return nil, err
	}
	return keyPair, nil
}

func GetRetryableTerraformErrors() map[string]string {
	retryableTerraformErrors := map[string]string{
		// The reason is unknown, but eventually these succeed after a few retries.
		".*unable to verify signature.*":                    "Failed due to transient network error.",
		".*unable to verify checksum.*":                     "Failed due to transient network error.",
		".*no provider exists with the given name.*":        "Failed due to transient network error.",
		".*registry service is unreachable.*":               "Failed due to transient network error.",
		".*connection reset by peer.*":                      "Failed due to transient network error.",
		".*TLS handshake timeout.*":                         "Failed due to transient network error.",
		".*Error: disassociating EC2 EIP.*does not exist.*": "Failed to delete EIP because interface is already gone",
		".*context deadline exceeded.*":                     "Failed due to kubernetes timeout, retrying.",
		".*http2: client connection lost.*":                 "Failed due to transient network error.",
	}
	return retryableTerraformErrors
}

func SetAcmeServer() string {
	acmeserver := os.Getenv("ACME_SERVER_URL")
	if acmeserver == "" {
		os.Setenv("ACME_SERVER_URL", "https://acme-staging-v02.api.letsencrypt.org/directory")
	}
	return acmeserver
}

func GetRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if region == "" {
		region = "us-west-2"
	}
	return region
}

func GetAwsAccessKey() string {
	key := os.Getenv("AWS_ACCESS_KEY_ID")
	if key == "" {
		key = "FAKE123-ABC"
	}
	return key
}

func GetAwsSecretKey() string {
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secret == "" {
		secret = "FAKE123-ABC"
	}
	return secret
}

func GetAwsSessionToken() string {
	return os.Getenv("AWS_SESSION_TOKEN")
}

func GetBuild() bool {
	if os.Getenv("SKIP_BUILD") == "true" {
		return false
	} else {
		return true
	}
}

func GetId() string {
	id := os.Getenv("IDENTIFIER")
	if id == "" {
		id = random.UniqueId()
	}
	id += "-" + random.UniqueId()
	return id
}

func CreateTestDirectories(t *testing.T, id string) error {
	gwd := g.GetRepoRoot(t)
	fwd, err := filepath.Abs(gwd)
	if err != nil {
		return err
	}
	tdd := fwd + "/test/data"
	err = os.Mkdir(tdd, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	tdd = fwd + "/test/data/" + id
	err = os.Mkdir(tdd, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	tdd = fwd + "/test/data/" + id + "/providers"
	err = os.Mkdir(tdd, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	tdd = fwd + "/test/data/" + id + "/plugins"
	err = os.Mkdir(tdd, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func Teardown(t *testing.T, directory string, options *terraform.Options, keyPair *aws.Ec2Keypair) {
	directoryExists := true
	_, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			directoryExists = false
		}
	}
	if directoryExists {
		_, err2 := terraform.DestroyE(t, options)
		if err2 != nil {
			// don't fail the test if destroying the cluster fails
			t.Logf("Error destroying cluster: %s", err2)
		}
		err := os.RemoveAll(directory)
		if err != nil {
			t.Logf("Failed to delete test data directory: %v", err)
		}
	}
	aws.DeleteEC2KeyPair(t, keyPair)
}
