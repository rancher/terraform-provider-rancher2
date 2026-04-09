## v4.1.0 (March 6, 2024)

FEATURES:

* Added support for `rancher2_pod_security_admission_configuration_template` resource and dataSource. See https://github.com/rancher/terraform-provider-rancher2/pull/1302

ENHANCEMENTS:

* [Docs] fix unit and add new example to the documentation for `rancher_cluster_v2`. See  https://github.com/rancher/terraform-provider-rancher2/pull/1309
* [Docs] add note linking to AWS docs for `eks_config_v2`. See https://github.com/rancher/terraform-provider-rancher2/pull/1247

BUGFIXES:


## [14.0.0](https://github.com/rancher/terraform-provider-rancher2/compare/v14.0.0...v14.0.0) (2026-04-09)


### Features

* Add new field `extended` to resource quota handlers ([#1811](https://github.com/rancher/terraform-provider-rancher2/issues/1811)) ([3f3da57](https://github.com/rancher/terraform-provider-rancher2/commit/3f3da57a5168c6b697b125d4d364774e375a2794))
* add reserved memory field ([#1510](https://github.com/rancher/terraform-provider-rancher2/issues/1510)) ([c1b4fa0](https://github.com/rancher/terraform-provider-rancher2/commit/c1b4fa0bbb74f2a9be9616419ede917504e8e5e0))
* added ipv6 suppport for eks operator ([#2005](https://github.com/rancher/terraform-provider-rancher2/issues/2005)) ([e6f04c7](https://github.com/rancher/terraform-provider-rancher2/commit/e6f04c7b0e4182b44f1222ee71c0117c22bcc0fb))
* cluster proxy config ([#1683](https://github.com/rancher/terraform-provider-rancher2/issues/1683)) ([e2a374a](https://github.com/rancher/terraform-provider-rancher2/commit/e2a374a01ab451ec90c6fd41f646b833309a8185))
* GitHub app provider ([#1695](https://github.com/rancher/terraform-provider-rancher2/issues/1695)) ([cf6e584](https://github.com/rancher/terraform-provider-rancher2/commit/cf6e5849de63cee6acc859b7bebceb61d345ce15))
* implement PC and PDB support for the fleet agent ([#1986](https://github.com/rancher/terraform-provider-rancher2/issues/1986)) ([ebf8d31](https://github.com/rancher/terraform-provider-rancher2/commit/ebf8d31fd977a53aa40730f5a8c147f1d5145993))
* implement support to nutanix machine config ([#2083](https://github.com/rancher/terraform-provider-rancher2/issues/2083)) ([a214fea](https://github.com/rancher/terraform-provider-rancher2/commit/a214fea018afd20222a5ece6058af1b8a30e3a04))
* keycloak OIDC auth provider ([#2033](https://github.com/rancher/terraform-provider-rancher2/issues/2033)) ([90d2e4e](https://github.com/rancher/terraform-provider-rancher2/commit/90d2e4e1f5afd278ff8730120bae297af8e3de26))
* support generic OIDC provider ([#1644](https://github.com/rancher/terraform-provider-rancher2/issues/1644)) ([4417468](https://github.com/rancher/terraform-provider-rancher2/commit/4417468b60abee4e969aa598564ed3d7bd3f4593))


### Bug Fixes

* [Docs] Add inherited_cluster_roles to docs ([#1342](https://github.com/rancher/terraform-provider-rancher2/issues/1342)) ([04d1088](https://github.com/rancher/terraform-provider-rancher2/commit/04d10888eb9e360f01c45a1dd22a1409516e1207))
* add all protocols to metadata ([#1794](https://github.com/rancher/terraform-provider-rancher2/issues/1794)) ([186d8f6](https://github.com/rancher/terraform-provider-rancher2/commit/186d8f6643c43911ea7b9ad98bb72fb7e080a5c8))
* add automation for release issues and tests ([#1692](https://github.com/rancher/terraform-provider-rancher2/issues/1692)) ([31d75e5](https://github.com/rancher/terraform-provider-rancher2/commit/31d75e5180721e1a451c50e7c7e5ae3e3f274ef6))
* add common environment variables to workflow ([#1742](https://github.com/rancher/terraform-provider-rancher2/issues/1742)) ([827540e](https://github.com/rancher/terraform-provider-rancher2/commit/827540e7e95aafe97cbbd890e714d86acdbf06a1))
* add config for release-please ([#1730](https://github.com/rancher/terraform-provider-rancher2/issues/1730)) ([575e79e](https://github.com/rancher/terraform-provider-rancher2/commit/575e79e8c9f7ee5b032ba571c8ece6bbb481f265))
* add deploy example of metrics operator ([#1608](https://github.com/rancher/terraform-provider-rancher2/issues/1608)) ([e28af8f](https://github.com/rancher/terraform-provider-rancher2/commit/e28af8ffcfed73a578e21dc42d3c1f09105bcc43))
* add FOSSA scanning workflow ([#1988](https://github.com/rancher/terraform-provider-rancher2/issues/1988)) ([7c865b3](https://github.com/rancher/terraform-provider-rancher2/commit/7c865b36803b553f3bb8586e60b90bbd15fb6f2c))
* add manual rc release and update permissions ([#1746](https://github.com/rancher/terraform-provider-rancher2/issues/1746)) ([24537b0](https://github.com/rancher/terraform-provider-rancher2/commit/24537b0fc7b986716d03d9b8771ad7c897b6e737))
* add permissions to update issues ([#1936](https://github.com/rancher/terraform-provider-rancher2/issues/1936)) ([a48759f](https://github.com/rancher/terraform-provider-rancher2/commit/a48759f52e5c531d65b7a7ca2828c9d9ce87c504))
* add search option exact match in principal ([#1331](https://github.com/rancher/terraform-provider-rancher2/issues/1331)) ([a3f15f1](https://github.com/rancher/terraform-provider-rancher2/commit/a3f15f14cca0446676e0d499f119798ccbd686ec))
* add support for AuthConfig Cognito ([#2029](https://github.com/rancher/terraform-provider-rancher2/issues/2029)) ([1f3fd79](https://github.com/rancher/terraform-provider-rancher2/commit/1f3fd793f7f86c6c96f26170227a143426f5286c))
* add support for data directories ([#1665](https://github.com/rancher/terraform-provider-rancher2/issues/1665)) ([637dbb9](https://github.com/rancher/terraform-provider-rancher2/commit/637dbb9755af488a279b7e72fb301c54d4bbf83d))
* adjust path and remove module context ([#1706](https://github.com/rancher/terraform-provider-rancher2/issues/1706)) ([e2c9c08](https://github.com/rancher/terraform-provider-rancher2/commit/e2c9c0822e377c28ee6bcd864bfdb947c9361b17))
* allow specifying sha ([#1764](https://github.com/rancher/terraform-provider-rancher2/issues/1764)) ([47e9f55](https://github.com/rancher/terraform-provider-rancher2/commit/47e9f5598e6cfb25d6b0e80c745c8a30d2469d43))
* break out the config pkg for better testing ([#1899](https://github.com/rancher/terraform-provider-rancher2/issues/1899)) ([80e2197](https://github.com/rancher/terraform-provider-rancher2/commit/80e21977be85183cc89fdd86990a5fa8a207d31d))
* build before testing ([#1903](https://github.com/rancher/terraform-provider-rancher2/issues/1903)) ([fed4427](https://github.com/rancher/terraform-provider-rancher2/commit/fed442767443ca63b01fe5dde75183ed6a8f3526))
* Bump golang.org/x/crypto from 0.33.0 to 0.35.0 ([#1530](https://github.com/rancher/terraform-provider-rancher2/issues/1530)) ([e617600](https://github.com/rancher/terraform-provider-rancher2/commit/e617600293904004211c5f7f6adb2f3ec1130ee3))
* change from "main" to "tracking" ([#1712](https://github.com/rancher/terraform-provider-rancher2/issues/1712)) ([35d9146](https://github.com/rancher/terraform-provider-rancher2/commit/35d9146a4ba995654c08bec2e6a4c70686000cc9))
* change IP for server URL in docs ([#1613](https://github.com/rancher/terraform-provider-rancher2/issues/1613)) ([338812e](https://github.com/rancher/terraform-provider-rancher2/commit/338812e8a4b170b86bb79c49970847e0a4f16687))
* change launch template version ([#1616](https://github.com/rancher/terraform-provider-rancher2/issues/1616)) ([8ccd8ec](https://github.com/rancher/terraform-provider-rancher2/commit/8ccd8ec600bd469ce055534a749fd16c5a75c731))
* change log from info to error ([#1604](https://github.com/rancher/terraform-provider-rancher2/issues/1604)) ([2e64ed4](https://github.com/rancher/terraform-provider-rancher2/commit/2e64ed4d4d69493b8924cc700b3c60c80a009b2d))
* check out repository ([#1751](https://github.com/rancher/terraform-provider-rancher2/issues/1751)) ([cdd8727](https://github.com/rancher/terraform-provider-rancher2/commit/cdd8727a4d4994ef75e8f33b531ebae563e9e086))
* clear space on runner before release ([#1848](https://github.com/rancher/terraform-provider-rancher2/issues/1848)) ([a5cb46c](https://github.com/rancher/terraform-provider-rancher2/commit/a5cb46c3eb83fa790d927f6ee04095b2710a74da))
* clear space on runners before release tests ([#1838](https://github.com/rancher/terraform-provider-rancher2/issues/1838)) ([e11d17d](https://github.com/rancher/terraform-provider-rancher2/commit/e11d17dfaa85658bf589d815c2ee5a646f3586e0))
* Correct genericoidc provider implementation and tests ([#1691](https://github.com/rancher/terraform-provider-rancher2/issues/1691)) ([0d7b485](https://github.com/rancher/terraform-provider-rancher2/commit/0d7b485f4f5289a9c7f098e618b685731fd3985d))
* do not check out branch that we tag ([#1932](https://github.com/rancher/terraform-provider-rancher2/issues/1932)) ([666f798](https://github.com/rancher/terraform-provider-rancher2/commit/666f798b3eab0cddb2dfd0b723a10ed7ced53584))
* do not use dot syntax in path ([#1705](https://github.com/rancher/terraform-provider-rancher2/issues/1705)) ([5821942](https://github.com/rancher/terraform-provider-rancher2/commit/5821942a237dd6f3c708fa21b3cf773e5b102df9))
* don't check user state on creation ([#1699](https://github.com/rancher/terraform-provider-rancher2/issues/1699)) ([4ebecdb](https://github.com/rancher/terraform-provider-rancher2/commit/4ebecdb4186a1676616da2fbcdfb29b1cfdca334))
* emit all response data for troubleshooting ([#1714](https://github.com/rancher/terraform-provider-rancher2/issues/1714)) ([d8490c3](https://github.com/rancher/terraform-provider-rancher2/commit/d8490c314b55ab63d0476a3ec74bc47012dd1bd2))
* enable manual rc release for older branches ([#1834](https://github.com/rancher/terraform-provider-rancher2/issues/1834)) ([afc0a9f](https://github.com/rancher/terraform-provider-rancher2/commit/afc0a9fd7118f9ba3b7c30c7d0b78817a0d2599b))
* ensure the correct workflow, users, and tags ([#1772](https://github.com/rancher/terraform-provider-rancher2/issues/1772)) ([e28d375](https://github.com/rancher/terraform-provider-rancher2/commit/e28d375ea4d4f5829cff4825d0844eef4a7e43c8))
* export GitHub scripts and validate in CI ([#1702](https://github.com/rancher/terraform-provider-rancher2/issues/1702)) ([4412d64](https://github.com/rancher/terraform-provider-rancher2/commit/4412d64f6113bc3a1bcabd9dd4dbd1dcd572ec97))
* fix script paths ([#1704](https://github.com/rancher/terraform-provider-rancher2/issues/1704)) ([0609f72](https://github.com/rancher/terraform-provider-rancher2/commit/0609f72d82257a82c66f19133f969c6b5163ceba))
* for manual releases skip the git validation ([#1788](https://github.com/rancher/terraform-provider-rancher2/issues/1788)) ([d67d532](https://github.com/rancher/terraform-provider-rancher2/commit/d67d532439a371b71ee8e505c0bad72430e18cae))
* force directory path ([#1710](https://github.com/rancher/terraform-provider-rancher2/issues/1710)) ([305ea61](https://github.com/rancher/terraform-provider-rancher2/commit/305ea616c608e710dec5db2fe8c0b22d38e54a38))
* github script uses slashes for comments ([#1863](https://github.com/rancher/terraform-provider-rancher2/issues/1863)) ([9a266f9](https://github.com/rancher/terraform-provider-rancher2/commit/9a266f9edcec8a18b6a1b0870eedfff1366b8a6b))
* goreleaser needs git state at the new tag ([#1768](https://github.com/rancher/terraform-provider-rancher2/issues/1768)) ([2d2805d](https://github.com/rancher/terraform-provider-rancher2/commit/2d2805dfa68c48043795f38d1b2616f21cc378aa))
* improve cleanup and skip in CI ([#1895](https://github.com/rancher/terraform-provider-rancher2/issues/1895)) ([8f5c657](https://github.com/rancher/terraform-provider-rancher2/commit/8f5c657973c1d18f193809a6791e39a06bc2b334))
* improve registry documentation ([#1664](https://github.com/rancher/terraform-provider-rancher2/issues/1664)) ([72bbd0c](https://github.com/rancher/terraform-provider-rancher2/commit/72bbd0c8ff703749c1ad584abf43d3524ccd10dc))
* improve test isolation logic ([#1879](https://github.com/rancher/terraform-provider-rancher2/issues/1879)) ([188dae9](https://github.com/rancher/terraform-provider-rancher2/commit/188dae9f0745cbe7609da762f47a1b4e480e6ae7))
* improve the documentation for cluster resources  ([#1805](https://github.com/rancher/terraform-provider-rancher2/issues/1805)) ([d2f8ef4](https://github.com/rancher/terraform-provider-rancher2/commit/d2f8ef4f65a7ff66eb83ce84020be4b11cdfcab0))
* inconsistent plan with EKS template ([#1639](https://github.com/rancher/terraform-provider-rancher2/issues/1639)) ([3140be6](https://github.com/rancher/terraform-provider-rancher2/commit/3140be67a25605893373c53fa9fef6ec849f34ea))
* login ([#2016](https://github.com/rancher/terraform-provider-rancher2/issues/2016)) ([edb74ae](https://github.com/rancher/terraform-provider-rancher2/commit/edb74ae0b807cbeb91b4da2ea059af8ed25d2d9d))
* make sure to only match the tag specified ([#1776](https://github.com/rancher/terraform-provider-rancher2/issues/1776)) ([cb60c2a](https://github.com/rancher/terraform-provider-rancher2/commit/cb60c2a78eaa0b44b1eb2d7d71356ec8b8c5b1d1))
* manual full release workflow ([#1759](https://github.com/rancher/terraform-provider-rancher2/issues/1759)) ([73f3e9c](https://github.com/rancher/terraform-provider-rancher2/commit/73f3e9c9f4e827310bba58f59dec59f103fe5be4))
* migrate login flow to /v1-public endpoint ([#1997](https://github.com/rancher/terraform-provider-rancher2/issues/1997)) ([3523df3](https://github.com/rancher/terraform-provider-rancher2/commit/3523df3ed8dceb1db465783a3170be6d1a1e0f33))
* missing quotes causing syntax error ([#1907](https://github.com/rancher/terraform-provider-rancher2/issues/1907)) ([6bfbfd4](https://github.com/rancher/terraform-provider-rancher2/commit/6bfbfd4fe6faee60c94f87cc9ffe288cc82c2591))
* need to quote json list ([#1875](https://github.com/rancher/terraform-provider-rancher2/issues/1875)) ([565c318](https://github.com/rancher/terraform-provider-rancher2/commit/565c3182b5e3d23231072a58b2e96afc4fa8d257))
* notifications to issues for manual rc ([#1923](https://github.com/rancher/terraform-provider-rancher2/issues/1923)) ([d05989f](https://github.com/rancher/terraform-provider-rancher2/commit/d05989fe7c0f4c791173c351768a0c92e185737a))
* prepare for v14 release ([#2012](https://github.com/rancher/terraform-provider-rancher2/issues/2012)) ([d396219](https://github.com/rancher/terraform-provider-rancher2/commit/d39621997494dbd17f0cd8b38fc8f4d03899d9c8))
* process imported config alone ([#1651](https://github.com/rancher/terraform-provider-rancher2/issues/1651)) ([bc8d237](https://github.com/rancher/terraform-provider-rancher2/commit/bc8d237f52edb1af8bbcfbb9a270c4d9ab2103eb))
* PSACT deletion hanging ([#1653](https://github.com/rancher/terraform-provider-rancher2/issues/1653)) ([1bccd68](https://github.com/rancher/terraform-provider-rancher2/commit/1bccd68900cd627fbe49545e15086ebe11fc601a))
* rc release calculator for new branches ([#2025](https://github.com/rancher/terraform-provider-rancher2/issues/2025)) ([cfad809](https://github.com/rancher/terraform-provider-rancher2/commit/cfad809bfda59b7ee8c6e7288b5280cb3550a5bd))
* remove broken resources ([#1693](https://github.com/rancher/terraform-provider-rancher2/issues/1693)) ([2785ae4](https://github.com/rancher/terraform-provider-rancher2/commit/2785ae4fc73abe826df791f7b02984dc8293ecf0))
* remove GKE enable logic ([#2000](https://github.com/rancher/terraform-provider-rancher2/issues/2000)) ([0b9e382](https://github.com/rancher/terraform-provider-rancher2/commit/0b9e382e29389c3996a43e346826bb9491f43d36))
* Remove ingress Nginx annotation examples ([#2038](https://github.com/rancher/terraform-provider-rancher2/issues/2038)) ([610b2f5](https://github.com/rancher/terraform-provider-rancher2/commit/610b2f56d3a6f4bd480e1605e5ff9275bfe18983))
* remove protected filter when getting branches in issue tracking workflow ([#2110](https://github.com/rancher/terraform-provider-rancher2/issues/2110)) ([b115402](https://github.com/rancher/terraform-provider-rancher2/commit/b1154021d93b9b5450c283859d5c4a8918dd2b26))
* remove pull target ([#2098](https://github.com/rancher/terraform-provider-rancher2/issues/2098)) ([3def343](https://github.com/rancher/terraform-provider-rancher2/commit/3def343e0c12e269b6fc1df2dcd6bd3acdc04b60))
* remove the default template ([#1780](https://github.com/rancher/terraform-provider-rancher2/issues/1780)) ([acc669c](https://github.com/rancher/terraform-provider-rancher2/commit/acc669cf0ad10c09d798d197d62aeb33ebc34208))
* repository owner is in github context ([#1891](https://github.com/rancher/terraform-provider-rancher2/issues/1891)) ([b23d521](https://github.com/rancher/terraform-provider-rancher2/commit/b23d521ee107dfebb9a04823a5d510e6d00a289a))
* run tests in different jobs ([#1859](https://github.com/rancher/terraform-provider-rancher2/issues/1859)) ([651634a](https://github.com/rancher/terraform-provider-rancher2/commit/651634a862d94779b02e0c2154a62a6cad47782d))
* security vulnerabilities and tracking workflow label ([#2130](https://github.com/rancher/terraform-provider-rancher2/issues/2130)) ([4e5d376](https://github.com/rancher/terraform-provider-rancher2/commit/4e5d3764d8702b1c66ddca1477ea1ba2691f3cf9))
* set bootstrap sha and release version ([#1755](https://github.com/rancher/terraform-provider-rancher2/issues/1755)) ([54b7dcf](https://github.com/rancher/terraform-provider-rancher2/commit/54b7dcf3dde80b52a1acaeef5c37b61a1547294b))
* set new release version ([#2021](https://github.com/rancher/terraform-provider-rancher2/issues/2021)) ([b466f31](https://github.com/rancher/terraform-provider-rancher2/commit/b466f31925f7b6d533e9e1ec8363646029d3df8a))
* short circuit the testing logic ([#1867](https://github.com/rancher/terraform-provider-rancher2/issues/1867)) ([42a638d](https://github.com/rancher/terraform-provider-rancher2/commit/42a638d5b012d17ff96ccadfb8b52900c034d57f))
* support exceptions in diff suppression ([#1808](https://github.com/rancher/terraform-provider-rancher2/issues/1808)) ([2bb8da1](https://github.com/rancher/terraform-provider-rancher2/commit/2bb8da1a6bf3fb23ffb6feea549eeda9f4b1506c))
* the first job should never skip ([#1887](https://github.com/rancher/terraform-provider-rancher2/issues/1887)) ([d6a3ea6](https://github.com/rancher/terraform-provider-rancher2/commit/d6a3ea6de5b7b2200b1a86d4cc4ea9cc51a7ae8a))
* typo in repository context ([#1883](https://github.com/rancher/terraform-provider-rancher2/issues/1883)) ([c5ed793](https://github.com/rancher/terraform-provider-rancher2/commit/c5ed7931de2049f7b5af8eded61a851d0200a35d))
* typo in script name ([#1915](https://github.com/rancher/terraform-provider-rancher2/issues/1915)) ([0257fdf](https://github.com/rancher/terraform-provider-rancher2/commit/0257fdf267cba828ace0667569393a6bc0d9336a))
* typo in variable name ([#1911](https://github.com/rancher/terraform-provider-rancher2/issues/1911)) ([694d6ee](https://github.com/rancher/terraform-provider-rancher2/commit/694d6ee981c52546119d4c6e912932dd95a0c501))
* update checkout ([#1842](https://github.com/rancher/terraform-provider-rancher2/issues/1842)) ([58979f1](https://github.com/rancher/terraform-provider-rancher2/commit/58979f1c7cd321719497c83ed0a535738070854a))
* Update compatibility matrix ([#1617](https://github.com/rancher/terraform-provider-rancher2/issues/1617)) ([1c0e84c](https://github.com/rancher/terraform-provider-rancher2/commit/1c0e84c49974f199aa12991afb823368ac33a7a1))
* update machine selector config example ([#1976](https://github.com/rancher/terraform-provider-rancher2/issues/1976)) ([bda0250](https://github.com/rancher/terraform-provider-rancher2/commit/bda025005fc5ae18427d8a7081db0d235ef5c558))
* update pre-release to see new tag format ([#1623](https://github.com/rancher/terraform-provider-rancher2/issues/1623)) ([bbbd78f](https://github.com/rancher/terraform-provider-rancher2/commit/bbbd78f9463bd8019904d8bd91cf9899a4297640))
* update rancher2 config map import command ([#1287](https://github.com/rancher/terraform-provider-rancher2/issues/1287)) ([1061ac6](https://github.com/rancher/terraform-provider-rancher2/commit/1061ac6a404a4dac078a815805e07e04c6fd726a))
* update release config to get proper version ([#1955](https://github.com/rancher/terraform-provider-rancher2/issues/1955)) ([c2e8d4a](https://github.com/rancher/terraform-provider-rancher2/commit/c2e8d4a5d6559076cb7e01a010377cd3fed14e76))
* update release manifest to force new version ([#1951](https://github.com/rancher/terraform-provider-rancher2/issues/1951)) ([c878c36](https://github.com/rancher/terraform-provider-rancher2/commit/c878c366a4b2a38dfc06854a5cc7df4a6021aae8))
* update release to 13.1.0 ([#1943](https://github.com/rancher/terraform-provider-rancher2/issues/1943)) ([a9ca6e7](https://github.com/rancher/terraform-provider-rancher2/commit/a9ca6e7f1fe55dfe8dd9470f1bb91783ba536efc))
* update release to 13.1.1 ([#1947](https://github.com/rancher/terraform-provider-rancher2/issues/1947)) ([b4df9f1](https://github.com/rancher/terraform-provider-rancher2/commit/b4df9f18173ad061a11bb2a23b2820533573e189))
* update release version ([#1968](https://github.com/rancher/terraform-provider-rancher2/issues/1968)) ([687f1a8](https://github.com/rancher/terraform-provider-rancher2/commit/687f1a814abac0b5feef7f3b108ed24dd36b6480))
* update versions in registry manifest ([#1963](https://github.com/rancher/terraform-provider-rancher2/issues/1963)) ([e75639c](https://github.com/rancher/terraform-provider-rancher2/commit/e75639cda6634091849d427105c9642ebbd18cba))
* use back tick instead of single quotes ([#1852](https://github.com/rancher/terraform-provider-rancher2/issues/1852)) ([b449a29](https://github.com/rancher/terraform-provider-rancher2/commit/b449a2950334b1455c9e4250be01761bca4a2051))
* use environment to pass data to the script ([#1919](https://github.com/rancher/terraform-provider-rancher2/issues/1919)) ([0b62b72](https://github.com/rancher/terraform-provider-rancher2/commit/0b62b72bde062cccfda622ab132c18415e4f858f))
* use modern import logic ([#1708](https://github.com/rancher/terraform-provider-rancher2/issues/1708)) ([1d0d125](https://github.com/rancher/terraform-provider-rancher2/commit/1d0d125ca93bd783b1c9cec9245a7e70288bafe9))
* use proper outputs for validating test ([#1871](https://github.com/rancher/terraform-provider-rancher2/issues/1871)) ([46c9209](https://github.com/rancher/terraform-provider-rancher2/commit/46c920972bc20271b9cfe2ed2cfdb8dc7dd78f44))
* use separate goreleaser files for manual ([#1784](https://github.com/rancher/terraform-provider-rancher2/issues/1784)) ([f9cc6bb](https://github.com/rancher/terraform-provider-rancher2/commit/f9cc6bb04f3217ccdd4377636f11c80bcb5412b3))
* v3 public fallback logic ([#2081](https://github.com/rancher/terraform-provider-rancher2/issues/2081)) ([f9c7a99](https://github.com/rancher/terraform-provider-rancher2/commit/f9c7a9910f2adb7cecbec5c3b3268e4ba4e4713a))

## v4.0.0 (February 5, 2024)

FEATURES:

* Added support for `machine_selector_files` on the `rancher2_cluster_v2` resource. See https://github.com/rancher/terraform-provider-rancher2/pull/1225
* Added support for `graceful_shutdown_timeout` for vSphere nodes. See https://github.com/rancher/terraform-provider-rancher2/pull/1228
* Added `inherited_cluster_roles` attribute to global_role. See https://github.com/rancher/terraform-provider-rancher2/pull/1242


ENHANCEMENTS:

* Changed default Ubuntu image for DigitalOcean from 16.04 -> 22.04. See  https://github.com/rancher/terraform-provider-rancher2/pull/1213
* Improved EKS node_group error message. See https://github.com/rancher/terraform-provider-rancher2/pull/1205
* Set `protect-kernel-defaults` on rancher_v2 clusters. See https://github.com/rancher/terraform-provider-rancher2/pull/1244
* Updated `default_system_registry` behavior at `rancher2_app_v2` resource. See https://github.com/rancher/terraform-provider-rancher2/pull/1265
* Updated RKE config machine pool schema and structure. See https://github.com/rancher/terraform-provider-rancher2/pull/1253
* Added `project_id` syntax for `rancher2_project_role_template_binding`. See https://github.com/rancher/terraform-provider-rancher2/pull/1272
* Add the `keypair_name` argument  to the `openstack_config` machine config and node template.  See https://github.com/rancher/terraform-provider-rancher2/pull/1235
* [Docs] updated the documentation for `rancher_cluster_v2`. See https://github.com/rancher/terraform-provider-rancher2/pull/1300
* [Docs] Remove `keypair_name` from `amazonec2_config`. See https://github.com/rancher/terraform-provider-rancher2/pull/942


BUGFIXES:

* Fixed `machine_selector_config` old schema to have the correct type. See https://github.com/rancher/terraform-provider-rancher2/pull/1223
* Allow machine pools to be created and scaled down to 0. See https://github.com/rancher/terraform-provider-rancher2/pull/1232
* Fixed the issue where the for_each in `rancher2_registry` keeps changing the order of the list that has not changed. See https://github.com/rancher/terraform-provider-rancher2/pull/1268
* Fixed double base64-encoding of `ca_bundle` field. See https://github.com/rancher/terraform-provider-rancher2/pull/1296
* [Docs] Fixed description for using `kms_key` `eks_config_v2`. See https://github.com/rancher/terraform-provider-rancher2/pull/892



## 3.1.1 (August 3, 2023)

FEATURES:



ENHANCEMENTS:

* [Docs] Add Terraform docs for cluster and fleet agent customization, PSACT support, and authentication ping `entity_field_id` for the v3.1.1 patch release. See [#1175](https://github.com/rancher/terraform-provider-rancher2/pull/1175)
* [Docs] Fix broken markdown in `rancher2_cluster` resource. See [#1180](https://github.com/rancher/terraform-provider-rancher2/pull/1180)
* [Docs] Update wording in registry resource. See [#1185](https://github.com/rancher/terraform-provider-rancher2/pull/1185)
* [Docs] Add example for multiple machine pools in RKE2. See [#957](https://github.com/rancher/terraform-provider-rancher2/pull/957)

BUG FIXES:



## 3.1.0 (June 25, 2023)

FEATURES:

* **New Resource** `rancher2_custom_user_token` - Provides configuration options to create, modify, and delete a user token. See [#932](https://github.com/rancher/terraform-provider-rancher2/pull/932) and [#1130](https://github.com/rancher/terraform-provider-rancher2/pull/1130)
* Add Cluster and Fleet Agent Deployment Customization support for users to customize the tolerations, affinity, and resources of a downstream agent.  See [#1137](https://github.com/rancher/terraform-provider-rancher2/pull/1137)
  * **New Argument** `cluster_agent_deployment_customization` - (Optional) Optional customization for cluster agent. For Rancher v2.7.5 and above (list)
  * **New Argument** `fleet_agent_deployment_customization` - (Optional) Optional customization for fleet agent. For Rancher v2.7.5 and above (list)
  * **New Argument** `append_tolerations` - (Optional) User defined tolerations to append to agent (list)
  * **New Argument** `override_affinity` - (Optional) User defined affinity to override default agent affinity (string)
  * **New Argument** `override_resource_requirements` - (Optional) User defined resource requirements to set on the agent (list)
* Add Pod Security Admission Configuration Template (PSACT) support with state migration logic for 1.25+ RKE and v2 prov clusters. See [#1119](https://github.com/rancher/terraform-provider-rancher2/pull/1119) and [#1117](https://github.com/rancher/terraform-provider-rancher2/pull/1117)
  * **New Argument** `default_pod_security_admission_configuration_template_name` - (Optional) Cluster default pod security admission configuration template name (string)
  * **New Argument** `default_pod_security_admission_configuration_template_name` - (Computed) Cluster V2 default pod security admission configuration template name (string)
  
ENHANCEMENTS:

* **New Argument** `entity_id_field` - (Optional) Entity ID for authentication config (string). See [#1163](https://github.com/rancher/terraform-provider-rancher2/pull/1163)
* [Docs] Add dev process and rc docs. See [#1138](https://github.com/rancher/terraform-provider-rancher2/pull/1138)
* Rancher machine hostname truncation. See [#1147](https://github.com/rancher/terraform-provider-rancher2/pull/1147)
* Refactor kubeconfig logic to use token from cached kubeconfig and replace invalid/expired tokens properly. [See #1158](https://github.com/rancher/terraform-provider-rancher2/pull/1158) and [#1165](https://github.com/rancher/terraform-provider-rancher2/pull/1165)
* Bump go-getter to 1.7.1. See [#1118](https://github.com/rancher/terraform-provider-rancher2/pull/1118)
* Bump google.golang.org/grpc to 1.53.0. See [#1167](https://github.com/rancher/terraform-provider-rancher2/pull/1167)

BUG FIXES:

* Add missing AKS node pool options. See [#1122](https://github.com/rancher/terraform-provider-rancher2/pull/1122)
* Verify `desired_size` in EKS node groups. See [#1126](https://github.com/rancher/terraform-provider-rancher2/pull/1126)
* Set DO userdata default empty for v2prov. See [#1121](https://github.com/rancher/terraform-provider-rancher2/pull/1121)
* Support old version HarvesterConfig. See [#1132](https://github.com/rancher/terraform-provider-rancher2/pull/1132)
* Fix 'unexpected end of JSON input' error when setting Pod Security Policy Template on new project. See [#1113](https://github.com/rancher/terraform-provider-rancher2/pull/1113)
* Fix Harvester `disk_size` default value. See [#1149](https://github.com/rancher/terraform-provider-rancher2/pull/1149)
* Consider all possible cluster states before passing them to StateChangeConf. See [#1114](https://github.com/rancher/terraform-provider-rancher2/pull/1114)

Your open source contributions are invaluable to us.

## 3.0.1 (June 7, 2023)

FEATURES:



ENHANCEMENTS:



BUG FIXES:

* Support old version HarvesterConfig. See [#1132](https://github.com/rancher/terraform-provider-rancher2/pull/1132)

## 3.0.0 (May 1, 2023)

FEATURES:

* Add support for new Azure features - node template tags, network interface, availability zone, and option to enable standard SKU. See [#1023](https://github.com/rancher/terraform-provider-rancher2/pull/1023) and [#1062](https://github.com/rancher/terraform-provider-rancher2/pull/1062)
* Add support for node group node role on EKS clusters. See [#1049](https://github.com/rancher/terraform-provider-rancher2/pull/1049)
* Allow setting vmAffinity for Harvester clusters. See [#1024](https://github.com/rancher/terraform-provider-rancher2/pull/1024) and [#1110](https://github.com/rancher/terraform-provider-rancher2/pull/1110)
* Update Harvester node driver to support multi nics and disks. See [#1051](https://github.com/rancher/terraform-provider-rancher2/pull/1051)

ENHANCEMENTS:

* [Docs] Update branching strategy and compatibility matrix. See [#1060](https://github.com/rancher/terraform-provider-rancher2/pull/1060)
* [Docs] Add example for authenticated cluster level private registry. See [#1055](https://github.com/rancher/terraform-provider-rancher2/pull/1055)
* [Docs] Remove macAddress in Harvester node driver. See [#1089](https://github.com/rancher/terraform-provider-rancher2/pull/1089)
* Add k8s 1.25 support | Remove logging cis v1 support. See [#1056](https://github.com/rancher/terraform-provider-rancher2/pull/1056)
* Add issue and pull request templates. See [#1063](https://github.com/rancher/terraform-provider-rancher2/pull/1063)
* Update Wrangler to v1.1.0. See [#1079](https://github.com/rancher/terraform-provider-rancher2/pull/1079)
* Bump golang.org/x/net from 0.2.0 to 0.7.0. See [#1078](https://github.com/rancher/terraform-provider-rancher2/pull/1078)
* Bump github.com/hashicorp/go-getter from 1.4.0 to 1.7.0. See [#1077](https://github.com/rancher/terraform-provider-rancher2/pull/1077)
* Add script to test Terraform RCs locally on darwin/unix and windows. See [#1082](https://github.com/rancher/terraform-provider-rancher2/pull/1082) and [#1085](https://github.com/rancher/terraform-provider-rancher2/pull/1085)
* Implement retry logic to enforce timeouts. See [#1033](https://github.com/rancher/terraform-provider-rancher2/pull/1033)

BUG FIXES:

* [Docs] Add machine_pool `annotations` property. See [#1041](https://github.com/rancher/terraform-provider-rancher2/pull/1041)
* Add new field for computed values in App v2 resource. See [#1021](https://github.com/rancher/terraform-provider-rancher2/pull/1021)
* Do not sort mirror endpoints. See [#1029](https://github.com/rancher/terraform-provider-rancher2/pull/1029)
* Update SDK and make machine pool `cloud_credential_secret_name optional` to fix plan bug. See [#1070](https://github.com/rancher/terraform-provider-rancher2/pull/1070)
* Fix intermittent 409 conflict when creating new `rancher2_project` with PSPTID. See [#1058](https://github.com/rancher/terraform-provider-rancher2/pull/1058)

## 1.25.0 (November 22, 2022)

FEATURES:



ENHANCEMENTS:



BUG FIXES:

* Update CHANGELOG with known bug info from #909. See [#1025](https://github.com/rancher/terraform-provider-rancher2/pull/1025)
* Revert "Add admission_configuration configuration." See [#1026](https://github.com/rancher/terraform-provider-rancher2/pull/1026)

## 1.24.2 (October 24, 2022)

FEATURES:

* Add Outscale support for node driver. See [#962](https://github.com/rancher/terraform-provider-rancher2/pull/962)
* Allow setting labels on nodes with RKE2. See [#951](https://github.com/rancher/terraform-provider-rancher2/pull/951)

ENHANCEMENTS:

* [Docs] Add a note about Azure AD auth configuration migration. See [#983](https://github.com/rancher/terraform-provider-rancher2/pull/983)
* [Docs] Remove tech preview reference for features that are already GA. See [#1000](https://github.com/rancher/terraform-provider-rancher2/pull/1000)
* Bump docker url version to 20.10. See [#970](https://github.com/rancher/terraform-provider-rancher2/pull/970)

BUG FIXES:

* Use existing cluster registration token if conflict. See [#997](https://github.com/rancher/terraform-provider-rancher2/pull/997)

KNOWN BUG:

* An update to the `admission_configuration` field within the kube-api service performed in [#909](https://github.com/rancher/terraform-provider-rancher2/pull/909) prevents provider upgrades from v1.24.1 to v1.24.2 in cases where `admission_configuration` was previously defined due to a type mismatch. See [#1011](https://github.com/rancher/terraform-provider-rancher2/issues/1011)

## 1.24.1 (September 1, 2022)

FEATURES:

* **New Argument:** `amazonec2_config.http_endpoint` - (Optional) Enables or disables the HTTP metadata endpoint on your instances (string). See [#944](https://github.com/rancher/terraform-provider-rancher2/pull/944)
* **New Argument:** `amazonec2_config.http_tokens` - (Optional) The state of token usage for your instance metadata requests (string). See [#944](https://github.com/rancher/terraform-provider-rancher2/pull/944)
* **New Argument:** `rke_config.aci_network_provider` - (Optional/Computed) ACI provider config for RKE network (list maxitems:63). See [#912](https://github.com/rancher/terraform-provider-rancher2/pull/912)

ENHANCEMENTS:

* Add test coverage for amazonec2 node template. See [#952](https://github.com/rancher/terraform-provider-rancher2/pull/952)
* Remove non working args from `amazonec2_config`. See [#960](https://github.com/rancher/terraform-provider-rancher2/pull/960)
* Add release checklist to README. See [#971](https://github.com/rancher/terraform-provider-rancher2/pull/971)
* Bump rancher/rancher and go client dependencies to support the ACI Network Provider. See [#959](https://github.com/rancher/terraform-provider-rancher2/pull/959)

BUG FIXES:

* Fix broken logo link in README. See [#924](https://github.com/rancher/terraform-provider-rancher2/pull/924)
* Fix resource creation error on RKE2 cluster for Azure [#876](https://github.com/rancher/terraform-provider-rancher2/pull/876)

## 1.24.0 (May 25, 2022)

FEATURES:

* Add Drain Before Delete timeout to cluster_v2 resource. See [#903](https://github.com/rancher/terraform-provider-rancher2/pull/903)
* Add k3s/rke2 etcd snapshot restore support for cluster v2. See [#920](https://github.com/rancher/terraform-provider-rancher2/pull/920)

BUG FIXES:

* Fixed RKE2 EC2 cluster creation with standard user token. See [#898](https://github.com/rancher/terraform-provider-rancher2/pull/898)
* Fixed incorrect delete timeout values for cluster resources. See [#906](https://github.com/rancher/terraform-provider-rancher2/pull/906)
* Fixed incorrect timeout variable used in the app delete function. See [#897](https://github.com/rancher/terraform-provider-rancher2/pull/897)

## 1.23.0 (March 31, 2022)

FEATURES:

* The Harvester node driver is now supported. See [#861](https://github.com/rancher/terraform-provider-rancher2/pull/861).
* Added support for KEv2 cluster certificate rotation. See [#882](https://github.com/rancher/terraform-provider-rancher2/pull/882).
* Added support for node pool health checks on KEv2 clusters. See [#889](https://github.com/rancher/terraform-provider-rancher2/pull/889).

ENHANCEMENTS:

* `hetzner_node_driver` now supports specifying `server_labels`. See [#851](https://github.com/rancher/terraform-provider-rancher2/pull/851).
* `drain_before_delete` can now be specified on `rancher2_cluster_v2` resources. See [#890](https://github.com/rancher/terraform-provider-rancher2/pull/890).
* Added support for `rancher2_cloud_credential` imported resources.
* Added `rancher2_node_template` label for server resources in `hetzner_config`. See [#657](https://github.com/rancher/terraform-provider-rancher2/issues/657).

BUG FIXES:

* Importing KEv2 clusters with the Rancher client have their config correctly rewritten. See [#36128](https://github.com/rancher/rancher/issues/36128).
* Fix race condition for empty kubeconfig generation for `rancher2_cluster_sync` resource. See [#849](https://github.com/rancher/terraform-provider-rancher2/issues/849).

## 1.22.2 (January 6, 2022)

FEATURES:



ENHANCEMENTS:

* Added k8s specialized verbs `bind`, `escalate`, `impersonate` and `use` support, to `rancher2_global_role` and `rancher2_role_template` rules

BUG FIXES:

* Fix `rancher2_cluster_sync.state_confirm` behaviour https://github.com/rancher/terraform-provider-rancher2/issues/797
* Fix `rancher2_cluster` monitoring flips on resource update https://github.com/rancher/terraform-provider-rancher2/issues/825
* Fix kube_config generation function at `rancher2_cluster`, `rancher2_cluster_v2` and `rancher2_cluster_sync` for Rancher 2.6.0 and above https://github.com/rancher/terraform-provider-rancher2/issues/789

## 1.22.1 (December 23, 2021)

FEATURES:



ENHANCEMENTS:


BUG FIXES:

* Fix `getClusterKubeconfig` function to properly generate kubeconfig when `rancher2_cluster.LocalClusterAuthEndpoint` is enabled for Rancher lower than v2.6.x https://github.com/rancher/terraform-provider-rancher2/issues/829

## 1.22.0 (December 22, 2021)

FEATURES:

* **New Argument:** `rancher2_cloud_credential.s3_credential_config` - (Optional) S3 config for the Cloud Credential. For Rancher 2.6.0 and above (list maxitems:1)
* **New Argument:** `rancher2_cluster.rke_config.enable_cri_dockerd` - (Optional) Enable/disable using cri-dockerd. Deafult: `false` (bool) https://github.com/rancher/terraform-provider-rancher2/issues/792
* **New Argument:** `rancher2_cluster.rke_config.private_registries.ecr_credential_plugin` - (Optional) ECR credential plugin config (list maxitems:1)
* **New Argument:** `rancher2_cluster_v2.local_auth_endpoint` - - (Optional) Cluster V2 local auth endpoint (list maxitems:1)
* **Deprecated Argument:** `rancher2_cluster_v2.rke_config.local_auth_endpoint` - (Deprecated) Use `rancher2_cluster_v2.local_auth` endpoint instead (list maxitems:1)

ENHANCEMENTS:

* Updated `rancher2_cluster_v2` docs adding labels and annotations arguments https://github.com/rancher/terraform-provider-rancher2/issues/784
* Updated `findClusterRegistrationToken` function checking for correct Cluster Registration Token https://github.com/rancher/terraform-provider-rancher2/issues/791
* Updated `getClusterKubeconfig` function to properly delete a cluster if cluster not available https://github.com/rancher/terraform-provider-rancher2/issues/788
* Updated `rancher2_machine_config_v2` resource to allow its use by Rancher standard users https://github.com/rancher/terraform-provider-rancher2/issues/824
* Updated `rancher2_cluster.eks_config_v2` argument to fix EKS launch template issue https://github.com/rancher/terraform-provider-rancher2/issues/820
* Updated go modules and acceptance tests to support rancher v2.6.3

BUG FIXES:

* Fix `rancher2_cluster_v2.rke_config.registries` sort to avoid false diff
* Fix `rancher2_machine_config_v2` to properly get updated

## 1.21.0 (October 18, 2021)

FEATURES:

* **New Data Source** `rancher2_principal` Use this data source to retrieve information about a Rancher v2 Principal resource
* **New Argument:** `rancher2_bootstrap.initial_password`  - (Optional/Computed/Sensitive) Initial password for Admin user. Default: `admin` (string)

ENHANCEMENTS:

* Added `IsConflict` function to retry on update v2 resources
* Added `RestartClients` function to restart Rancher2 clients
* Updated `CatalogV2Client` function to also retry if err `IsNotFound` or `IsForbidden`
* Refactored `activateNodeDriver` and `activateKontainerDriver` functions to retry if got conflict on node driver action
* Refactored v2 resources CRUD functions to be defined at resources files. Added retry if got conflict on update
* Updated `rancher2_cluster_v2` to support creation using Rancher2 standard user token
* Refactored `config.NormalizeURL` function to return error
* Updated `rancher2_bootstrap.current_password` argument to be Computed/Sensitive. Please be sure to remove this field from tf file before update.

BUG FIXES:

* Fixed `waitAllCatalogV2Downloaded` function avoiding race condition

## 1.20.1 (October 5, 2021)

FEATURES:

* **New Data Source:** `rancher2_config_map_v2` - Provides Rancher configMap v2 data source. Available at Rancher v2.5.x and above.
* **New Resource:** `rancher2_config_map_v2` - Provides Rancher configMap v2 resource. Available at Rancher v2.5.x and above.

ENHANCEMENTS:

* Updated go modules and acceptance tests to support rancher v2.6.1
* Updated `waitForRancherLocalActive` function to allow `rancher2_bootstrap` works when using Rancher [restricted-admin](https://rancher.com/docs/rancher/v2.6/en/admin-settings/rbac/global-permissions/#restricted-admin) at Rancher 2.6.x
* Updated `rancher2_cluster.aks_config_v2` schema and structure to fix aks cluster import errors https://github.com/rancher/terraform-provider-rancher2/issues/757 https://github.com/rancher/terraform-provider-rancher2/issues/771

BUG FIXES:

* Fixed `expandClusterEKSConfigV2` function to avoid provider crash https://github.com/rancher/terraform-provider-rancher2/issues/753
* Fixed `rancher2_cluster` resource update to properly update eks v2 and gke v2 clusters

## 1.20.0 (September 17, 2021)

FEATURES:

* **New Argument:** `rancher2_cluster.aks_config_v2` - (Optional) The Azure AKS v2 configuration for creating/import `aks` Clusters. Conflicts with `aks_config`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` `k3s_config` and `rke_config`. For Rancher v2.6.0 and above (list maxitems:1)
* **New Argument:** `rancher2_cloud_credential.azure_credential_config.environment` - (Optional/Computed) Azure environment (e.g. AzurePublicCloud, AzureChinaCloud) (string)
* **New Argument:** `rancher2_cloud_credential.azure_credential_config.tenant_id` - (Optional/Computed) Azure Tenant ID (string)
* **New Attribute:** `rancher2_cluster.cluster_registration_token.insecure_node_command` - (Computed) Insecure node command to execute in a imported k8s cluster (string)
* **New Attribute:** `rancher2_cluster.cluster_registration_token.insecure_windows_node_command` - (Computed) Insecure windows command to execute in a imported k8s cluster (string)
* **New Attribute:** `rancher2_cloud_credential.amazonec2_credential_config.default_region` - (Optional) AWS default region (string)
* **New Resource:** `rancher2_machine_config_v2` - Provides a Rancher v2 Machine config v2 resource. Available as tech preview at Rancher v2.6.0 and above.
* **New Resource:** `rancher2_cluster_v2` - Provides Rancher cluster v2 resource to manage RKE2 and K3S cluster. Available as tech preview at Rancher v2.6.0 and above.
* **New Data Source:** `rancher2_cluster_v2` - Provides Rancher cluster v2 resource to manage RKE2 and K3S cluster. Available at Rancher v2.6.0 and above.

ENHANCEMENTS:

* Updated go modules and acceptance tests to support rancher v2.6.0
* Updated `rancher2_cluster.rke_config` schema to support rancher v2.6.0 https://github.com/rancher/rke/pull/2409
* Updated `rancher2_cluster.gke_config_v2` schema to support rancher v2.6.0 https://github.com/rancher/gke-operator/pull/49
* Updated `rancher2_cluster.eks_config_v2` schema to support rancher v2.6.0 https://github.com/rancher/eks-operator/pull/38
* Updated `rancher2_cluster.gke_config_v2` schema to support rancher v2.6.0 https://github.com/rancher/rancher/issues/34291
* Updated docs, adding note to use `rancher2_bootstrap` resource on Rancher v2.6.0 and above

BUG FIXES:

* Updated `rancher2_project_role_template_binding` `rancher2_cluster_role_template_binding` resources, setting `user_id` and `group_id` arguments as computed
* Updated `rancher2_cluster.aks_config_v2` to:
  * disable default value for `node_pools.max_count` and `node_pools.min_count` https://github.com/rancher/rancher/issues/34752
  * set optional arguments as computed for imported clusters https://github.com/rancher/rancher/issues/34758
* Updated `InfoAppV2` function to proper escape url query params  https://github.com/rancher/terraform-provider-rancher2/issues/739

## 1.17.2 (August 25, 2021)

FEATURES:



ENHANCEMENTS:

* Added verb `own` to policy rule
* Updated `WaitForClusterState` function to check for condition last update before return error

BUG FIXES:



## 1.17.1 (August 18, 2021)

FEATURES:

* **New Argument:** `rancher2_cluster.fleet_workspace_name` - (Optional/Computed) Fleet workspace name (string)

ENHANCEMENTS:



BUG FIXES:

* Fix `rancher2_cluster` resource update to not reset fleet workspace name
* Fix `rancher2_node_template` resource to proper update `cloud_credential_id` and `use_internal_ip_address` arguments

## 1.17.0 (August 12, 2021)

FEATURES:

* **New Argument:** `rancher2_cluster.rke_config.ingress.tolerations` - (Optional) Ingress add-on tolerations (list)
* **New Argument:** `rancher2_cluster.rke_config.monitoring.tolerations` - (Optional) Monitoring add-on tolerations (list)
* **New Argument:** `rancher2_cluster.rke_config.network.tolerations` - (Optional) Network add-on tolerations (list)
* **New Argument:** `rancher2_cluster.rke_config.dns.options` - (Optional) DNS add-on options (map)
* **New Argument:** `rancher2_cluster.rke_config.dns.tolerations` - (Optional) DNS add-on tolerations (list)
* **New Argument:** `rancher2_cluster.oke_config.enable_private_control_plane` - (Optional) Specifies whether Kubernetes API endpoint is a private IP only accessible from within the VCN. Default `false` For Rancher v2.5.10 and above (bool)
* **New Data Source:** `rancher2_storage_class_v2` - Provides Rancher Storage Class v2 data source. Available at Rancher v2.5.x and above.
* **New Resource:** `rancher2_storage_class_v2` - Provides Rancher Storage Class v2 resource. Available at Rancher v2.5.x and above.

ENHANCEMENTS:

* Added `tolerations` schema, structure and tests
* Updated `rancher2_cluster` resource to properly generate cluster registration token
* Minor `rancher2_catalog_v2` and `rancher2_secret_v2` datasource docs update
* Added verb `deletecollection` to policy rule
* Updated `WaitForClusterState` function to check for cluster transitioning before return error

BUG FIXES:

* Updated `rancher2_notifier` resource to be replaced on update
* Fixed `rancher2_cluster.eks_config_v2` to avoid false diff
* Updated `rancher2_notifier` resource to be replaced on update
* Updated `rancher2_cluster` docs to proper format yaml examples

## 1.16.0 (July 15, 2021)

FEATURES:

* **New Argument:** `rancher2_auth_config_keycloak.entity_id` - (Optional/Computed) KeyCloak Client ID field (string)
* **New Argument:** `rancher2_auth_config_activedirectory.start_tls` - (Optional/Computed) Enable start TLS connection (bool)
* **New Argument:** `rancher2_node_pool.drain_before_delete` - (Optional) Drain nodes before delete (bool)

ENHANCEMENTS:

* Added timeout error message to `CatalogV2Client`, `getObjectV2ByID` and `GetCatalogV2List` functions
* Updated `rancher2_bootstrap` resource to wait until `local` cluster is active
* Updated `rancher2_cluster.rke_config.cloud_provider.name` argument from `Optional/Computed` to `Optional`
* Updated `rancher2_cluster` resource to replace RKE cluster API info instead of update, if `rancher2_cluster.rke_config` has been updated
* Updated `rancher2_project` resource to replace project API info instead of update
* Updated `rancher2_node_template.engine_install_url` argument to be `computed`
* Updated Rancher to v2.5.9
* Updated golang to v1.16.5 and added darwin arm64 build

BUG FIXES:

* Fixed `rancher2_cluster.gke_config_v2.cluster_addons` to be optional

## 1.15.1 (May 21, 2021)

FEATURES:



ENHANCEMENTS:



BUG FIXES:

* Added timeout to `CatalogV2Client` function when getting new catalog v2 client

## 1.15.0 (May 20, 2021)

FEATURES:

* **Deprecated Argument:** `rancher2_cluster.aks_config.tag` - (Deprecated) Use `tags` argument instead as []string
* **New Argument:** `rancher2_cluster.aks_config.tags` - (Optional/Computed) Tags for Kubernetes cluster. For example, `["foo=bar","bar=foo"]` (list)
* **New Argument:** `rancher2_cluster.agent_env_vars` - (Optional) Optional Agent Env Vars for Rancher agent. For Rancher v2.5.6 and above (list)
* **Deprecated provider Argument:** `retries` - (Deprecated) Use timeout instead
* **New provider Argument:** `timeout` - (Optional) Timeout duration to retry for Rancher connectivity and resource operations. Default: `"120s"`
* **New Argument:** `rancher2_cluster.oke_config.pod_cidr` - (Optional) A CIDR IP range from which to assign Kubernetes Pod IPs (string)
* **New Argument:** `rancher2_cluster.oke_config.service_cidr` - (Optional) A CIDR IP range from which to assign Kubernetes Service IPs (string)

ENHANCEMENTS:

* Added timeout to `CatalogV2Client` function when getting new catalog v2 client

BUG FIXES:

* Fixed `rancher2_cluster.hetzner_config.UsePrivateNetwork` with proper json field name

## 1.14.0 (May 7, 2021)

FEATURES:

* **New Argument:** `rancher2_cluster.oke_config.limit_node_count` - (Optional) The maximum number of worker nodes. Can limit `quantity_per_subnet`. Default `0` (no limit) (int)
* **New Argument:** `rancher2_cluster.rke_config.ingress.default_backend` - (Optional) Enable ingress default backend. Default: `true` (bool)
* **New Argument:** `rancher2_cluster.rke_config.ingress.http_port` - (Optional/Computed) HTTP port for RKE Ingress (int)
* **New Argument:** `rancher2_cluster.rke_config.ingress.https_port` - (Optional/Computed) HTTPS port for RKE Ingress (int)
* **New Argument:** `rancher2_cluster.rke_config.ingress.network_mode` - (Optional/Computed) Network mode for RKE Ingress (string)
* **New Argument:** `rancher2_cluster.rke_config.ingress.update_strategy` - (Optional) RKE ingress update strategy (list Maxitems: 1)
* **New Argument:** `rancher2_cluster.rke2_config` - (Optional/Computed) The RKE2 configuration for `rke2` Clusters. Conflicts with `aks_config`, `eks_config`, `gke_config`, `oke_config`, `k3s_config` and `rke_config` (list maxitems:1)
* **New Argument:** `rancher2_cluster_sync.wait_alerting` - (Optional) Wait until alerting is up and running. Default: `false` (bool)
* **New Argument:** `rancher2_cluster.gke_config_v2` - (Optional) The Google GKE V2 configuration for `gke` Clusters. Conflicts with `aks_config`, `eks_config`, `eks_config_v2`, `gke_config`, `oke_config`, `k3s_config` and `rke_config`. For Rancher v2.5.8 and above (list maxitems:1)
* **New Argument:** `rancher2_cloud_credential.google_credential_config` - (Optional) Google config for the Cloud Credential (list maxitems:1)

ENHANCEMENTS:

* Updated `rancher2_catalog_v2` schema resource, defining conflict between `git_repo` and `url` arguments
* Improved `rancher2_cluster_sync` with new cluster state check method and new option to wait until alerting is enabled
* Updated go mod to support Rancher `v2.5.8`
* Updated acceptance tests to use Rancher `v2.5.8`

BUG FIXES:

* Fix `rancher2_node_pool` resource, adding `forcenew` property to not updatable arguments
* Fix `rancher2_cluster` resource, fixing provider crash if `cluster_monitoring_input` argument is deleted
* Fix `rancher2_project` resource, fixing provider crash if `project_monitoring_input` argument is deleted
* Fix `rancher2_catalog_v2` resource, just setting default `git_branch` value if `git_repo` is specified
* Fix `rancher2_cluster.eks_config_v2` argument, setting `private_access`, `public_access` and `secrets_encryption` as computed argument, removing default value

## 1.13.0 (March 31, 2021)

FEATURES:

* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.image_id` - (Optional) The EKS node group image ID (string)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.launch_template` - (Optional) The EKS node groups launch template (list Maxitem: 1)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.launch_template.id` - (Required) The EKS node group launch template ID (string)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.launch_template.name` - (Optional/Computed) The EKS node group launch template name (string)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.launch_template.version` - (Optional) The EKS node group launch template version. Default: `1` (int)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.request_spot_instances` - (Optional) Enable EKS node group request spot instances (bool)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.resource_tags` - (Optional) The EKS node group resource tags (map)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.spot_instance_types` - (Optional) The EKS node group sport instace types (list string)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.subnets` - (Optional) The EKS node group subnets (list string)
* **New Argument:** `rancher2_cluster.eks_config_v2.node_groups.user_data` - (Optional) The EKS node group user data (string)
* **New Argument:** `rancher2_cluster_sync.wait_catalogs` - (Optional) Wait until all catalogs are downloaded and active. Default: `false` (bool)
* **New Attribute:** `rancher2_cluster.eks_config_v2.node_groups.version` - (Computed) The EKS node group version (string)
* **New Attribute:** `rancher2_app_v2.system_default_registry` - (Computed) The system default registry of the app (string)
* **New Data Source:** `rancher2_secret_v2` - Provides a Rancher V2 Secret V2 data source
* **New Resource:** `rancher2_secret_v2` - Provides a Rancher V2 Secret V2 resource

ENHANCEMENTS:

* Updated go mod to support Rancher `v2.5.7`
* Updated acceptance tests to use Rancher `v2.5.7`
* Updated `rancher2_cluster_sync` to allow wait until all catalogs are downloaded and active

BUG FIXES:

* Fix `rancher2_app_v2` to respect Rancher system default registry
* Fix `rancher2_cluster.eks_config_v2` to deploy properly EKS clusters
* Fix `rancher2_catalog_v2` to wait until `downloaded` status

## 1.12.0 (March 05, 2021)

FEATURES:

* **New Argument:** `rancher2_node_template.node_taints` - (Optional) Node taints. For Rancher v2.3.3 and above (List)
* **New Argument:** `rancher2_cluster.aks_config.load_balancer_sku` - (Optional/Computed) Load balancer type (basic | standard). Must be standard for auto-scaling
* **New Argument:** `rancher2_cluster.rke_config.services.etc.backup_config.timeout` - (Optional/Computed) Set timeout in seconds for etcd backup. For Rancher v2.5.6 and above
* **New Data Source:** `rancher2_global_role` - Provides a Rancher V2 Global Role data source
* **New Resource:** `rancher2_global_role` - Provides a Rancher V2 Global Role resource
* **New Resource:** `rancher2_feature` - Provides a Rancher V2 Feature resource. For Rancher v2.5.0 and above

ENHANCEMENTS:

* Updated `rancher2_node_template.openstack_config` to support `boot_from_volume` and related arguments
* Added `password` as valid `cluster_template_questions` type to `rancher2_cluster` resource
* Preserve `cluster_template_answers` for `cluster_template_questions` of type `password` in `rancher2_cluster` resource to avoid misleading diffs
* Added `nodes` attribute reference to `rancher2_cluster_sync` resource
* Updated go mod to support Rancher `v2.5.6`
* Updated acceptance tests to use Rancher `v2.5.6`
* Added retry to get k8s default version, if getting forbidden or server error
* Added retry to get V2 catalogs and apps, if getting server error

BUG FIXES:

* Fixed cluster and project resource for update monitoring version properly
* Fixed `rancher2_app_v2` resource, added retry to GetAppV2OperationByID if got apierr 500
* Fixed `rancher2_cluster` docs, annotations and labels argument description

## 1.11.0 (January 08, 2021)

FEATURES:

* **New Argument:** `rancher2_node_template.hetzner_config` - (Optional) Hetzner config for the Node Template (list maxitems:1)
* **New Argument:** `rancher2_cluster.rke_config.dns.linear_autoscaler_params` - (Optional) LinearAutoScalerParams dns config (list Maxitem: 1)
* **New Argument:** `rancher2_cluster.rke_config.dns.update_strategy` - (Optional) DNS update strategy (list Maxitems: 1)
* **New Argument:** `rancher2_notifier.dingtalk_config` - (Optional) Dingtalk config for notifier (list maxitems:1)
* **New Argument:** `rancher2_notifier.msteams_config` - (Optional) MSTeams config for notifier (list maxitems:1)
* **New Data Source:** `rancher2_global_dns_provider` - Provides a Rancher V2 Global DNS Provider data source
* **New Resource:** `rancher2_global_dns` - Provides a Rancher V2 Global DNS resource
* **New Resource:** `rancher2_global_dns_provider` - Provides a Rancher V2 Global DNS Provider resource

ENHANCEMENTS:

* Updated `rancher2_app_v2.chart_version` as optional/computed argument. Deploying latest app v2 version if `chart_version` is not provided
* Updated `rancher2_app_v2.wait` default value to `true`
* Updated go mod to support Rancher `v2.5.4`
* Updated acceptance tests to use Rancher `v2.5.4`

BUG FIXES:

* Fixed `rancher2_cluster` resource, added retry when enabling cluster monitoring and got apierr 500. https://github.com/rancher/rancher/issues/30188
* Fixed `rancher2_cluster` datasource error, when `rke_config.services.kube_api.secrets_encryption_config.custom_config` or `rke_config.services.kube_api.event_rate_limit.configuration` are set. https://github.com/rancher/terraform-provider-rancher2/issues/546
* Fixed `rancher2_cluster_template` required argument definition on docs
* Fixed `Apps & marketplace` guide for Rancher v2.5.0 format
* Fixed doc examples for activedirectory, freeipa and openldap auth providers
* Fixed `rancher2_app_v2` resource to properly pass global values to sub charts. https://github.com/rancher/terraform-provider-rancher2/issues/545
* Fixed `rancher2_app_v2` resource to don't override name nor namespace on App v2 not certified by rancher
* Fixed `rancher2_cluster` docs, adding missed `gke_config.enable_master_authorized_network` argument

## 1.10.6 (November 11, 2020)

FEATURES:



ENHANCEMENTS:


BUG FIXES:

* Fixed `flattenClusterTemplateRevisions` func to avoid crash on `rancher2_cluster_template` resource at some circumstances

## 1.10.5 (November 11, 2020)

FEATURES:

* **Deprecated Argument:** `rancher2_cluster.eks_import` - (Optional) Use `rancher2_cluster.eks_config_v2` instead. For Rancher v2.5.0 and above
* **New Argument:** `rancher2_cluster.eks_config_v2` - (Optional) EKS cluster import and new management support. For Rancher v2.5.0 and above

ENHANCEMENTS:

* Updated go mod to support Rancher `v2.5.2`
* Updated acceptance tests to use Rancher `v2.5.2`
* Improved `rancher2_bootstrap` on resource creation. Number of retires on `bootstrapDoLogin` function can be configured with `retries` provider argument
* Updated `rancher2_catalog_v2` contextualized resource id with `cluster_id` prefix
* Updated `rancher2_app_v2` contextualized resource id with `cluster_id` prefix
* Updated `rancher2_app_v2` to show helm operation log if fail
* Updated `rancher2_app_v2.values` argument as sensitive

BUG FIXES:

* Fixed `rancher2_cluster.rke_config.upgrade_strategy.drain` argument to set false value properly
* Fixed `Apps & marketplace` guide for Rancher v2.5.0 format
* Fixed `rancher2_app_v2.values` argument to avoid false diff
* Fixed `rancher2_cluster_role_template_binding` and  `rancher2_cluster_role_template_binding` arguments to forceNew on update

## 1.10.4 (October 29, 2020)

FEATURES:

* **New Argument:** `rancher2_cluster.oke_config` - (Optional) Oracle OKE configuration
* **New Argument:** `rancher2_node_template.openstack_config.application_credential_id` - (Optional) OpenStack application credential id
* **New Argument:** `rancher2_node_template.openstack_config.application_credential_name` - (Optional) OpenStack application credential name
* **New Argument:** `rancher2_node_template.openstack_config.application_credential_secret` - (Optional) OpenStack application credential secret
* **New Argument:** `rancher2_notifier.dingtal_config` - (Optional) Dingtalk config for notifier. For Rancher v2.4.0 and above (list maxitems:1)
* **New Argument:** `rancher2_notifier.msteams_config` - (Optional) MSTeams config for notifier. For Rancher v2.4.0 and above (list maxitems:1)
* **New Argument:** `rancher2_cluster.eks_import` - (Optional) EKS cluster import and new management support. For Rancher v2.5.0 and above
* **New Argument:** `rancher2_bootstrap.ui_default_landing` - (Optional) Set default ui landing on Rancher bootstrap. For Rancher v2.5.0 and above
* **New Data Source:** `rancher2_catalog_v2` - Support new Rancher catalog V2 datasource. For Rancher v2.5.0 and above
* **New Resource:** `rancher2_catalog_v2` - Support new Rancher catalog V2 resource. For Rancher v2.5.0 and above
* **New Resource:** `rancher2_app_v2` - Support new Rancher app V2 resource. For Rancher v2.5.0 and above

ENHANCEMENTS:

* Added new computed `ca_cert` argument at `rancher2_cluster` resource and datasource
* Delete `rancher2_app` if created and got timeout to be active
* Updated golang to v1.14.9 and removing vendor folder
* Updated go mod to support Rancher `v2.5.1`
* Added dingtal_config and msteams_config arguments at rancher2_notifier resource. go code and docs
* Improved `rancher2_cluster_sync` wait for cluster monitoring
* Improved `rancher2_bootstrap` on resource creation. `bootstrapDoLogin` function will retry 3 times user/pass login before fail
* Updated acceptance tests to use Rancher `v2.5.1`, k3s `v1.18.9-k3s1` and cert-manager `v1.0.1`
* Added new `Apps & marketplace` guide for Rancher v2.5.0

BUG FIXES:

* Fix `rke_config.monitoring.replicas` argument to set default value to 1 if monitoring enabled
* Fix Rancher auth config apply on activedirectory, freeipa and openldap providers
* Fix `rancher2_cluster.rke_config.upgrade_strategy.drain` argument to set false value properly


## 1.10.3 (September 14, 2020)

FEATURES:



ENHANCEMENTS:



BUG FIXES:

* Fix `Error: string is required` upgrading rancher2 provider from v1.10.0 or lower

## 1.10.2 (September 10, 2020)

FEATURES:



ENHANCEMENTS:

* Updated go mod, vendor files and provider tests to support rancher 2.4.8 and k3s v1.18.8-k3s1
* Added `rancher2_cluster_sync.state_confirm` argument to wait until active status is confirmed a number of times
* Added `syslog_config.enable_tls` argument to cluster and project logging

BUG FIXES:

* Fix `rke_config.cloud_provider.name` argument to not be validated
* Fix `rancher2_certificate` resource update
* Fix false diff if `rancher2_project.project_monitoring_input` not specified
* Fix `rancher2_token.ttl` argument to work properly on Rancher up to v2.4.7
* Fix `rancher2_namespace.resource_quota` argument to computed
* Fix `rancher2_app` resource to wait until created/updated

## 1.10.1 (August 27, 2020)

FEATURES:



ENHANCEMENTS:

* Added `nsg` support on `azure_config` argument on `rancher2_node_template` resource
* Updated go mod, vendor files and provider tests to support rancher 2.4.6
* Added aws kms key id support to `rancher2_node_template`

BUG FIXES:

* Fix `rke_config.event_rate_limit.configuration` argument to work properly
* Fix cluster and project role template binding doc files name
* Fix `rancher2_cluster_sync` resource error if referred cluster deleted out of band
* Fix `rancher2_namespace` and `rancher2_project` resources error if destroyed by not global admin user
* Fix `rancher2_app` resource error if referred project deleted out of band
* Fix `rancher2_app` doc typo on `target_namespace` argument description
* Fix `rancher2_cluster` and `rancher2_project` resources error if created with monitoring enabled by not global admin user
* Fix `rancher2_token` to set annotations and labels as computed attibutes
* Fix `rke_config.secrets_encryption_config.custom_config` argument to work properly
* Fix `rancher2_token.ttl` argument to work properly on Rancher v2.4.6
* Fix `rancher2_project` resource applying `pod_security_policy_template_id` argument on creation

## 1.10.0 (July 29, 2020)

FEATURES:

* **Deprecated Argument:** `rancher2_cluster.enable_cluster_istio` - Deploy istio using `rancher2_app` resource instead
* **New Argument:** `rancher2_cluster.istio_enabled` - (Computed) Is istio enabled at cluster?

ENHANCEMENTS:

* Added `wait` argument to rancher2_app
* Added configurable retry logic when Rancher responds with "405 method not allowed" for `rancher2_node_template` resource
* Added drone pipeline definition to publish provider at terraform registry
* Updated docs to terraform registry format

BUG FIXES:

* Fixes on `rancher2_cluster_template` resource:
  * Update default revision. Related to https://github.com/rancher/terraform-provider-rancher2/issues/393
  * Import. Related to https://github.com/rancher/terraform-provider-rancher2/issues/386
  * Delete old template revisions. Related to https://github.com/rancher/terraform-provider-rancher2/issues/397
* Fixed import resource description on doc files
* Fixed bootstrap link on doc website

## 1.9.0 (June 29, 2020)

FEATURES:



ENHANCEMENTS:

* Updated acceptance tests:
  * run rancher HA environment on k3s v1.18.2-k3s1
  * integrated rancher update scenario from v2.3.6 to v2.4.5
* Updated local cluster on `rancher2_bootstrap` resource, due to issue https://github.com/rancher/rancher/issues/16213
* Added `load_balancer_sku` argument to `azure_cloud_provider` configuration
* Added `nodelocal` argument to `rke_config.dns` argument on `rancher2_cluster` resource
* Added `view` verb to `rules` argument for `rancher2_node_template` resource
* Updated golang to v1.13, modules and vendor files
* Updated Rancher support to v2.4.5
* Added full feature to `rke_config.monitoring` argument
* Added `external` as allowed value on `rke_config.cloud_provider` argument on `rancher2_cluster` resource
* Added `region` argument on `gke_config` for `rancher2_cluster` resource
* Updated `annotations` and `labels` arguments to supress diff when name contains `cattle.io/` or `rancher.io/`

BUG FIXES:

* Fixed `nodeTemplateStateRefreshFunc` function on `rancher2_node_template` resource to check if returned error is forbidden
* Updated `rancher2_app` resource to fix local cluster scoped catalogs
* Updated api bool fields with default=true to `*bool`. Related to https://github.com/rancher/types/pull/1083
* Fixed update on `rancher2_cluster_template` resource. Related to https://github.com/terraform-providers/terraform-provider-rancher2/issues/365

## 1.8.3 (April 09, 2020)

FEATURES:



ENHANCEMENTS:



BUG FIXES:

* Fix project alert group and alert rule datasource and resoruce documentation
* Added `version` argument to `cluster_monitoring_input` argument on `rancher2_cluster` and `rancher2_project` resources
* Fixed rancher timeout on bootstrapping

## 1.8.2 (April 02, 2020)

FEATURES:



ENHANCEMENTS:

* Added `fixNodeTemplateID` to fix `rancher2_node_template` ID upgrading up to v2.3.3. Issue [#195](https://github.com/terraform-providers/terraform-provider-rancher2/issues/195)
* Updated rnacher to v2.4.2 on acceptance test

BUG FIXES:

* Fix `upgrading` state on resourceRancher2ClusterUpdate() function
* Updated process for getting rke supported kubernetes version
* Set `version` argument on `rancher2_catalog` as computed

## 1.8.1 (March 31, 2020)

FEATURES:



ENHANCEMENTS:



BUG FIXES:

* Fix init provider if api_url is dependent of infra that is not yet deployed

## 1.8.0 (March 31, 2020)

FEATURES:

* **New Data Source:** `rancher2_cluster_scan`

ENHANCEMENTS:

* Added `wait_monitoring` argument to `rancher2_cluster_sync` resource
* Added `retries` config argument and `isRancherActive()` function
* Updated go modules and vendor files to rancher v2.4.0
* Updated rancher to v2.4.0 and k3s to v1.17.4-k3s1 on acceptance tests
* New rancher v2.4.0 features:
  * Added `group_principal_id` argument to `rancher2_global_role_binding` resource
  * Added `k3s_config` argument to `rancher2_cluster` datasource and resource
  * Added `version` argument to `rancher2_catalog` datasource and resource
  * Added `upgrade_strategy` argument to `rke_config` on `rancher2_cluster` resource
  * Added `scheduled_cluster_scan` argument on `rancher2_cluster` and `rancher2_cluster_template` resources
  * Added `rancher2_cluster_scan` datasource
* Added `fixNodeTemplateID` to fix `rancher2_node_template` ID upgrading up to v2.3.3. Issue [#195](https://github.com/terraform-providers/terraform-provider-rancher2/issues/195)

BUG FIXES:

* Added `enable_json_parsing` argument to cluster and project logging
* Sync resource delete with rancher API
* Fix recipient update on cluster and project alert groups

## 1.7.3 (March 24, 2020)

FEATURES:

* **New Data Source:** `rancher2_pod_security_policy_template`
* **New Resource:** `rancher2_pod_security_policy_template`

ENHANCEMENTS:

* Updated `rancher/norman` go modules and vendor files
* Added `plugin` optional value `none` to `rke_config` argument on `rancher2_cluster` resource
* Updated multiline arguments to trim spaces by default and avoid false diff
* Updated `rancher/types` go modules and vendor files
* Added `mtu` argument to `rke_config.network` argument on `rancher2_cluster` resource
* Added `custom_target_config` argument to `rancher2_cluster_logging` and `rancher2_project_logging` resources
* Updated `aks_config`, `eks_config` and `gke_config` arguments ti proper updte `rancher2_cluster` resource

BUG FIXES:

* Fix `audit_log.configuration.policy` argument to `rke_config.services.kube_api` argument on `rancher2_cluster` resource
* Added `plugin` optional value `none` to `rke_config` argument on `rancher2_cluster` resource
* Updated multiline arguments to trim spaces by default and avoid false diff
* Updated `private_key_file` definition for openstack driver on `rancher2_node_template` docs
* Updated `private_key_file` definition for openstack driver on `rancher2_node_template` docs
* Fixed `rke_config.cloud_provider.aws_cloud_provider.global` argument as computed to avoid false diff

## 1.7.2 (January 28, 2020)

FEATURES:



ENHANCEMENTS:

* Added `refresh` argument to `rancher2_catalog` resource
* Added `name` and `is_external` arguments to `rancher2_user` datasource
* Added `delete_not_ready_after_secs` and `node_taints` arguments to `node_pool` resource
* Added `delete_not_ready_after_secs` and `node_taints` arguments to `rancher2_node_pool` resource
* Updated `github.com/rancher/types` and `github.com/rancher/norman` go modules and vendor files to support rancher v2.3.3
* Splitted schema, structure and test `cluster_rke_config_services` files for every rke service
* Added `ssh_cert_path` argument to `rke_config` argument on `rancher2_cluster` resource
* Added `audit_log`, `event_rate_limit` and `secrets_encryption_config` arguments to `rke_config.services.kube_api` argument on `rancher2_cluster` resource
* Added `generate_serving_certificate` argument to `rke_config.services.kubelet` argument on `rancher2_cluster` resource
* Added `driver_id` argument on `rancher2_node_template` resource to reference user created `rancher2_node_driver`

BUG FIXES:

* Fix `template_revisions` update on `rancher2_cluster_template` resource
* Fix `rke_config.services.kube_api.policy` argument on `rancher2_cluster` resource
* Fix `data` argument set as sensitive on `rancher2_secret` resource

## 1.7.1 (December 04, 2019)

FEATURES:



ENHANCEMENTS:

* Added GetRancherVersion function to provider config
* Updated `vsphere_config` argument schema on `rancher2_node_template` resource to support Rancher v2.3.3 features
* Updated rancher to v2.3.3 and k3s to v0.10.2 on acceptance tests

BUG FIXES:

* Set `annotations` argument as computed on `rancher2_node_template` resource
* Added `rancher2_node_template` resource workaround on docs when upgrade Rancher to v2.3.3

## 1.7.0 (November 20, 2019)

FEATURES:

* **New Resource:** `rancher2_token`

ENHANCEMENTS:

* Added `always_pull_images` argument on `kube_api` argument on `rke_config` argument for `rancher2_clusters` resource
* Added resource deletion if not getting active state on creation for `rancher2_catalog` resource
* Updated rancher to v2.3.2 and k3s to v0.10.1 on acceptance tests
* Added `desired nodes` support on `eks_config` argument on `rancher2_cluster` resource
* Added `managed disk` support on `azure_config` argument on `rancher2_node_template` resource
* Migrated provider to use `terraform-plugin-sdk`
* Updated `rancher2_etcd_backup` documentation

BUG FIXES:

* Fix `password` argument update for `rancher2_catalog` resource
* Fix `rancher2_app` update issue on Rancher v2.3.2
* Fix: set `key` argument as sensitive on `rancher2_certificate` resource.
* Fix continuous diff issues on `rancher2_project` resource
* Fix `pod_security_policy_template_id` update on `rancher2_project` resource
* Fix continuous diff issues on `rancher2_namespace` resource

## 1.6.0 (October 08, 2019)

FEATURES:

* **New Data Source:** `rancher2_cluster_alert_group`
* **New Data Source:** `rancher2_cluster_alert_rule`
* **New Data Source:** `rancher2_cluster_template`
* **New Data Source:** `rancher2_notifier`
* **New Data Source:** `rancher2_project_alert_group`
* **New Data Source:** `rancher2_project_alert_rule`
* **New Data Source:** `rancher2_role_template`
* **New Resource:** `rancher2_auth_config_keycloak`
* **New Resource:** `rancher2_auth_config_okta`
* **New Resource:** `rancher2_cluster_alert_group`
* **New Resource:** `rancher2_cluster_alert_rule`
* **New Resource:** `rancher2_cluster_sync`
* **New Resource:** `rancher2_cluster_template`
* **New Resource:** `rancher2_notifier`
* **New Resource:** `rancher2_project_alert_group`
* **New Resource:** `rancher2_project_alert_rule`
* **New Resource:** `rancher2_role_template`

ENHANCEMENTS:

* Added `monitoring_input` argument to define monitoring config for `rancher2_cluster` and `rancher2_project`
* Improved capitalization/spelling/grammar/etc in docs

BUG FIXES:

* Fix `expandAppExternalID` function on `rancher2_app` resource. Function was generating a wrong `ExternalID` catalog URL, on `cluster` and `project` scope
* Fix `flattenMultiClusterApp` function on `rancher2_multi-cluster_app` resource. Function wasn't updating fine `catalog_name`, `template_name` and/or `template_version` arguments, when contains char `-`
* Fix: set `value_yaml` multiline argument as base64 encoded
* Fix: removed `restricted` and `unrestricted` values checking for `default_pod_security_policy_template_id` argument on `rancher2_cluster` resource

## 1.5.0 (September 06, 2019)

FEATURES:

* **New Data Source:** `rancher2_app`
* **New Data Source:** `rancher2_certificate`
* **New Data Source:** `rancher2_multi_cluster_app`
* **New Data Source:** `rancher2_node_template`
* **New Data Source:** `rancher2_secret`
* **New Resource:** `rancher2_certificate`
* **New Resource:** `rancher2_app`
* **New Resource:** `rancher2_multi_cluster_app`
* **New Resource:** `rancher2_secret`

ENHANCEMENTS:

* Updated default image to `canonical:UbuntuServer:18.04-LTS:latest` on Azure node template
* Added `folder` argument on `s3_backup_config`
* Updated `github.com/rancher/types` and `github.com/rancher/norman` go modules and vendor files to support rancher v2.2.8
* Updated rancher to v2.2.8 and k3s to v0.8.0 on acceptance tests
* Added `key_pair_name` argument on `eks_config` argument on `rancher2_cluster` resource
* Set `kubernetes_version` argument as required on `eks_config` argument on `rancher2_cluster` resource
* Set `quantity` argument as optional with default value `1` on `rancher2_node_pool` resource. Added validation that value >= 1

BUG FIXES:

* Fix: `container_resource_limit` argument update issue on `rancher2_namespace` and `rancher2_project` resources update
* Fix: `sidebar_current` definition on datasources docs
* Fix: set `access_key` and `secret_key` arguments as optional on `s3_backup_config`
* Fix: crash `rancher2_cluster`  datasource and resource if `enableNetworkPolicy` doesn't exist
* Fix: don't delete builtin cluster nor node drivers from rancher on tf destroy
* Fix: wrong updates on not changed sensitive arguments on `rancher2_cluster_logging` and `rancher2_project_logging` resources

## 1.4.1 (August 16, 2019)

FEATURES:

ENHANCEMENTS:

BUG FIXES:

* Fix: auth issue when using `access_key` and `secret_key`

## 1.4.0 (August 15, 2019)

FEATURES:

* **New Data Source:** `rancher2_catalog`
* **New Data Source:** `rancher2_cloud_credential`
* **New Data Source:** `rancher2_cluster`
* **New Data Source:** `rancher2_cluster_driver`
* **New Data Source:** `rancher2_cluster_logging`
* **New Data Source:** `rancher2_cluster_role_template_binding`
* **New Data Source:** `rancher2_etcd_backup`
* **New Data Source:** `rancher2_global_role_binding`
* **New Data Source:** `rancher2_namespace`
* **New Data Source:** `rancher2_node_driver`
* **New Data Source:** `rancher2_node_pool`
* **New Data Source:** `rancher2_project_logging`
* **New Data Source:** `rancher2_project_role_template_binding`
* **New Data Source:** `rancher2_registry`
* **New Data Source:** `rancher2_user`
* **New Resource:** `rancher2_global_role_binding`
* **New Resource:** `rancher2_registry`
* **New Resource:** `rancher2_user`

ENHANCEMENTS:

* Set `session_token` argument as sensitive on `eks_config` argument on `rancher2_cluster` resource
* Added `wait_for_cluster` argument on `rancher2_namespace` and `rancher2_project` resources
* Set default value to `engine_install_url` argument on `rancher2_node_template` resource
* Added `enable_cluster_monitoring` argument to `rancher2_cluster` resource and datasource
* Added `enable_project_monitoring` argument to `rancher2_project` resource and datasource
* Added `token` argument on `cluster_registration_token` argument to rancher2_cluster resource and datasource
* Set default value to `engine_install_url` argument on `rancher2_node_template` resource
* Added `custom_ca` argument on etcd `s3_backup_config` on `rancher2_cluster` and `rancher2_etcd_backup` resources
* Updated `github.com/rancher/types` and `github.com/rancher/norman` go modules and vendor files to support rancher v2.2.6
* Updated rancher to v2.2.6 and k3s to v0.7.0 on acceptance tests
* Added cluster and project scope support on `rancher2_catalog` resource and datasource
* Updated `provider` config validation to enable bootstrap and resource creation at same run
* Added `container_resource_limit` argument on `rancher2_namespace` and `rancher2_project` resources and datasources
* Added `pod_security_policy_template_id` on `rancher2_project` resource

BUG FIXES:

* Fix: `toArrayString` and `toMapString` functions to check `nil` values
* Fix: Set `kubernetes_version` argument as required on `aks_config` argument on `rancher2_cluster` resource
* Fix: Set `security_groups`, `service_role`, `subnets` and `virtual_network` arguments as optional to `eks_config` argument on `rancher2_cluster` resource
* Fix: Removed `docker_version` argument from `rancher2_node_template` resource

## 1.3.0 (June 26, 2019)

FEATURES:

ENHANCEMENTS:

* Added `scheduler` argument to `services`-`rke_config` argument on `rancher2_cluster` resource

BUG FIXES:

* Fix: index out of range issue on `vsphere_cloud_provider`-`cloud_provider`-`rke_config` argument on `rancher2_cluster` resource

## 1.2.0 (June 12, 2019)

FEATURES:

* **New Data Source:** `rancher2_project`

ENHANCEMENTS:

* Added `cluster_auth_endpoint` argument to `rancher2_cluster` resource
* Added `default_pod_security_policy_template_id` argument to `rancher2_cluster` resource
* Added `enable_network_policy` argument to `rancher2_cluster` resource
* Updated acceptance tests
  * k3s version updated to v0.5.0
  * Rancher version updated to v2.2.4

BUG FIXES:

* Fix: set default value to `true` on `ignore_docker_version`-`rke_config` argument on `rancher2_cluster` resource
* Fix: set default value to `false` on `pod_security_policy`-`services`-`rke_config` argument on `rancher2_cluster` resource
* Fix: typo on `boot2docker_url`-`vsphere_config` argument name on `rancher2_node_template` resource docs
* Fix: set `monitor_delay` and `monitor_timeout` fields as string type for openstack load_balancer config on `cloud_provider`-`rke_config` argument on `rancher2_cluster` resource
* Fix: Updated `rancher2_etcd_backup` resource to work on rancher v2.2.4

## 1.1.0 (May 29, 2019)

FEATURES:

ENHANCEMENTS:

* Added `default_project_id` & `system_project_id` attributes to `rancher2_cluster` resource
* Added support to move `rancher2_namespace` resource to a rancher project when import
* Added support to terraform 0.12

BUG FIXES:

* Fix: Updated `flattenNamespace` function on `rancher2_namespace` resource to avoid no empty plan if `resource_quota` is not specified
* Fix: Updated `rke_config` argument for openstack cloud_provider on `rancher2_cluster` resource:
  * Removed `used_id` field on global argument in favour of `username` following [k8s openstack cloud provider docs](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/provider-configuration.md#global-required-parameters)
  * Set computed=true on optional field to avoid no empty plan if not specified

## 1.0.0 (May 14, 2019)

* Initial Terraform Ecosystem Release


## v0.2.0-rc5 (Unreleased)

FEATURES:

ENHANCEMENTS:

* Updated `rancher2_cluster` `rke_config` argument to support `aws_cloud_provider` config
* Updated k3s version to v0.4.0 to run acceptance tests
* Added support to openstack and vsphere drivers on `rancher2_cloud_credential` resource
* Added support to openstack and vsphere drivers on `rancher2_node_template` resource

BUG FIXES:

* Fix: Updated `rancher2_cluster` resource to save correctly S3 and cloud providers passwords on `rke_config`
* Updated `rancher2_cloud_credential` resource to save correctly S3 password
* Updated `rancher2_etcd_backup` resource to save correctly S3 password

## v0.2.0-rc4 (Unreleased)

FEATURES:

* **New Resource:** `rancher2_bootstrap`
* **New Resource:** `rancher2_cloud_credential`
* **New Resource:** `rancher2_cluster_driver`
* **New Resource:** `rancher2_etcd_backup`

ENHANCEMENTS:

* Added `.drone.yml` file to also support run rancher pipeline
* Added `rancher2_node_pool` resource tests
* Added `rancher2_auth_config_*` resource tests
* Updated and reviewed docs format
* Added support to rancher v2.2.x
* Updated `rancher2_cluster` `rke_config` argument to support:
  * etcd service `backup_config` with local and S3 storage backends
  * `dns` config
  * `weave` network provider
* Splitted resources into own schema, structure and import files.
* Added support to amazonec2, azure and digitalocean drivers on `rancher2_cloud_credential` resource
* Added support to local and S3 storage backends on `rancher2_etcd_backup` resource

BUG FIXES:

* Fix: drone build image to golang:1.12.3 to fix go fmt issues
* Fix: removed test on apply for `rancher2_auth_config_*` resources
* Fix: updated `api_url` field as required on provider.go
* Fix: updated `rancher2_namespace` move to a project after import it from k8s cluster

## v0.2.0-rc3 (Unreleased)

FEATURES:

ENHANCEMENTS:

* Added `Sensitive: true` option to fields with sensible data

BUG FIXES:

* Fix: set rke cluster `cloud_provider_vsphere` disk and network as optional and computed fields

## v0.2.0-rc2 (Unreleased)

FEATURES:

ENHANCEMENTS:

* Added `Sensitive: true` option to fields with sensible data
* Added `kube_config` computed field on cluster resources
* Added `ami` and `associate_worker_node_public_ip` fields for `eks_config` on cluster resources
* Added all available fields for rke_config on cluster resources
* Added `manifest_url` and `windows_node_command` fields for `cluster_registration_token` on cluster resources
* Added `creation` argument on `etcd` service for rke_config on cluster resources

BUG FIXES:

* Fix: added updating pending state on cluster resource update
* Fix: checking if `cluster_registration_token` exists on cluster resource creation
* Fix: typo on `gke_config` credential field on cluster resource
* Fix: Updated auth resources to avoid permission denied error

## 0.1.0-rc1 (Unreleased)

FEATURES:

* **New Resource:** `rancher2_auth_config_activedirectory`
* **New Resource:** `rancher2_auth_config_adfs`
* **New Resource:** `rancher2_auth_config_azuread`
* **New Resource:** `rancher2_auth_config_freeipa`
* **New Resource:** `rancher2_auth_config_github`
* **New Resource:** `rancher2_auth_config_openldap`
* **New Resource:** `rancher2_auth_config_ping`
* **New Resource:** `rancher2_catalog`
* **New Resource:** `rancher2_cluster`
* **New Resource:** `rancher2_cluster_logging`
* **New Resource:** `rancher2_cluster_role_template_binding`
* **New Resource:** `rancher2_namespace`
* **New Resource:** `rancher2_node_driver`
* **New Resource:** `rancher2_node_pool`
* **New Resource:** `rancher2_node_template`
* **New Resource:** `rancher2_project`
* **New Resource:** `rancher2_project_logging`
* **New Resource:** `rancher2_project_role_template_binding`
* **New Resource:** `rancher2_setting`

ENHANCEMENTS:

* First release candidate of the rancher2 provider.
* resource/rancher2_cluster: support for providers:
  * Amazon EKS
  * Azure AKS
  * Google GKE
  * Imported
  * RKE
    * Cloud providers adding node pools
    * Custom
* resource/rancher2_cluster_logging: support for providers:
  * Elasticsearch
  * Fluentd
  * Kafka
  * Splunk
  * Syslog
* resource/rancher2_namespace: quota limits support on Rancher v2.1.x or higher
  * Amazon EC2
  * Azure
  * Digitalocean
* resource/rancher2_project: quota limits support on Rancher v2.1.x or higher
* resource/rancher2_project_logging: support for providers:
  * Elasticsearch
  * Fluentd
  * Kafka
  * Splunk
  * Syslog
* resource/rancher2_node_template: support for providers:

BUG FIXES:
