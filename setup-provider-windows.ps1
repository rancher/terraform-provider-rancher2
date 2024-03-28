<#
This Windows script is to be used to test Terraform RCs locally. It downloads
and installs the given RC.

./setup-provider-windows.ps1 <provider> <version>

Example

./setup-provider-windows.ps1 rancher v3.0.0-rc1

If you get a cert revocation check error when running the script, curl could
not identify a valid certificate to download the zip. Add --ssl-no-revoke or
--insecure to the curl command and rerun. Alternatives such as white listing
the download link may also work.
#>

$PROVIDER=$args[0]
$VERSION=$args[1]
$VERSION_TAG=$VERSION.TrimStart("v")

$DIR="%APPDATA%\terraform.d\plugins\terraform.local\local\${PROVIDER}\${VERSION_TAG}\windows_amd64"
new-item $DIR -itemtype directory

write-host "Downloading zip..."
$URL="https://github.com/rancher/terraform-provider-${PROVIDER}/releases/download/${VERSION}/terraform-provider-${PROVIDER}_${VERSION_TAG}_windows_amd64.zip"
echo $URL
curl.exe -LO $URL -o terraform-provider-${PROVIDER}_${VERSION_TAG}_windows_amd64.zip

write-host "Expanding zip..."
Expand-Archive -Force -Path terraform-provider-${PROVIDER}_${VERSION_TAG}_windows_amd64.zip -DestinationPath ${DIR}/terraform-provider-${PROVIDER}

write-host "Cleaning up..."
Remove-Item -Force -Path terraform-provider-${PROVIDER}_${VERSION_TAG}_windows_amd64.zip

write-host "Terraform is ready to test!"
