# Changelog

## [14.0.0](https://github.com/rancher/terraform-provider-rancher2/compare/v14.0.0...v14.0.0) (2026-02-26)


### Features

* implement PC and PDB support for the fleet agent ([#2035](https://github.com/rancher/terraform-provider-rancher2/issues/2035)) ([db0748a](https://github.com/rancher/terraform-provider-rancher2/commit/db0748aad4266ee0533e2d4ec88ab23a05fdf173))


### Bug Fixes

* [Docs] Add inherited_cluster_roles to docs ([#2056](https://github.com/rancher/terraform-provider-rancher2/issues/2056)) ([451adf4](https://github.com/rancher/terraform-provider-rancher2/commit/451adf403afe9f5ecc3ccb349890381eba8527f1))
* add FOSSA scanning workflow ([#2032](https://github.com/rancher/terraform-provider-rancher2/issues/2032)) ([78fe9ad](https://github.com/rancher/terraform-provider-rancher2/commit/78fe9ad355410a4c97d0390840e53c9ad1d96665))
* add search exact match option for principal ([#2067](https://github.com/rancher/terraform-provider-rancher2/issues/2067)) ([cb0cfc7](https://github.com/rancher/terraform-provider-rancher2/commit/cb0cfc722fe00a5a507d15464ec14107ec403282))
* Bump golang.org/x/crypto from 0.33.0 to 0.35.0 ([#2059](https://github.com/rancher/terraform-provider-rancher2/issues/2059)) ([5203465](https://github.com/rancher/terraform-provider-rancher2/commit/52034653925cafd44c8fbc924a6a0b52d0a82574))
* login ([#2019](https://github.com/rancher/terraform-provider-rancher2/issues/2019)) ([3a41053](https://github.com/rancher/terraform-provider-rancher2/commit/3a41053af986308c6665b7b86dc53cdde195f9ca))
* prepare for v14 release ([#2015](https://github.com/rancher/terraform-provider-rancher2/issues/2015)) ([8ca2e99](https://github.com/rancher/terraform-provider-rancher2/commit/8ca2e997f339c883dad32c20d4ed20b656e3c614))
* rc release calculator for new branches ([#2028](https://github.com/rancher/terraform-provider-rancher2/issues/2028)) ([5b566d7](https://github.com/rancher/terraform-provider-rancher2/commit/5b566d752446a2d6888afad5a7353588e8d8b9cc))
* Remove ingress Nginx annotation examples ([#2077](https://github.com/rancher/terraform-provider-rancher2/issues/2077)) ([5917b45](https://github.com/rancher/terraform-provider-rancher2/commit/5917b4523a9dbcc0211454e3fe68a4e9a6c7ad3e))
* set new release version ([#2024](https://github.com/rancher/terraform-provider-rancher2/issues/2024)) ([cffd050](https://github.com/rancher/terraform-provider-rancher2/commit/cffd050fea400699e4b06cfbf12de3ae86f32f6b))
* update machine selector config example ([#2076](https://github.com/rancher/terraform-provider-rancher2/issues/2076)) ([70cf584](https://github.com/rancher/terraform-provider-rancher2/commit/70cf584fdeef697309d6cec7c171214bbecf0372))
* update rancher2 config map import command ([#2070](https://github.com/rancher/terraform-provider-rancher2/issues/2070)) ([1020115](https://github.com/rancher/terraform-provider-rancher2/commit/1020115b1770c420dc2288d07d6c07ffcf9c22cc))
* v3 public fallback logic ([#2086](https://github.com/rancher/terraform-provider-rancher2/issues/2086)) ([fa05d99](https://github.com/rancher/terraform-provider-rancher2/commit/fa05d99c7feb3d2a2b5d76acd40c20c487909707))

## [13.1.4](https://github.com/rancher/terraform-provider-rancher2/compare/v13.1.3...v13.1.4) (2025-12-19)


### Bug Fixes

* update release version ([#1971](https://github.com/rancher/terraform-provider-rancher2/issues/1971)) ([6d16fee](https://github.com/rancher/terraform-provider-rancher2/commit/6d16feeca6870d0cbb781b10a0fd5ade94f7ba96))
* update versions in registry manifest ([#1966](https://github.com/rancher/terraform-provider-rancher2/issues/1966)) ([71ed6ea](https://github.com/rancher/terraform-provider-rancher2/commit/71ed6ea4f8526a992a0cd2e4e85e15debdb23d7d))

## [13.1.3](https://github.com/rancher/terraform-provider-rancher2/compare/v13.1.2...v13.1.3) (2025-12-16)


### Features

* add reserved memory field ([#1817](https://github.com/rancher/terraform-provider-rancher2/issues/1817)) ([47d79b9](https://github.com/rancher/terraform-provider-rancher2/commit/47d79b9620b9569d084487db076984f84cad3257))
* Add support for machine configurations ([6a2e4e6](https://github.com/rancher/terraform-provider-rancher2/commit/6a2e4e658c3079a5dcfe96bdffea764d93581ee3))
* GitHub app provider ([#1928](https://github.com/rancher/terraform-provider-rancher2/issues/1928)) ([f296c23](https://github.com/rancher/terraform-provider-rancher2/commit/f296c23c0597dee9656b893328989dde9806b308))
* support generic OIDC provider ([#1644](https://github.com/rancher/terraform-provider-rancher2/issues/1644)) ([#1648](https://github.com/rancher/terraform-provider-rancher2/issues/1648)) ([a97b5f5](https://github.com/rancher/terraform-provider-rancher2/commit/a97b5f539b77742fe6ac7d9d63e8c83a5bc3b67d))


### Bug Fixes

* [v9] change launch template version ([#1621](https://github.com/rancher/terraform-provider-rancher2/issues/1621)) ([97e8f04](https://github.com/rancher/terraform-provider-rancher2/commit/97e8f048da4b1667cb68617cbb5814b4cf707ec0))
* add all protocols to metadata ([#1797](https://github.com/rancher/terraform-provider-rancher2/issues/1797)) ([f1cf58e](https://github.com/rancher/terraform-provider-rancher2/commit/f1cf58e222b3364e6381bf608d1192070f371346))
* add common environment variables to workflow ([#1745](https://github.com/rancher/terraform-provider-rancher2/issues/1745)) ([788d9eb](https://github.com/rancher/terraform-provider-rancher2/commit/788d9eb6c363a8963bc2f92171668003c7dcc1a5))
* add config for release-please ([#1740](https://github.com/rancher/terraform-provider-rancher2/issues/1740)) ([3adbc38](https://github.com/rancher/terraform-provider-rancher2/commit/3adbc3862e503e0b9cf9fee11863efb936bbd112))
* add deploy example of metrics operator ([#1608](https://github.com/rancher/terraform-provider-rancher2/issues/1608)) ([e28af8f](https://github.com/rancher/terraform-provider-rancher2/commit/e28af8ffcfed73a578e21dc42d3c1f09105bcc43))
* Add Doc for Rancher pod security admission configuration template ([752ed1b](https://github.com/rancher/terraform-provider-rancher2/commit/752ed1b56532cc4656894ec2c5eb0918261558d5))
* add manual rc release and update permissions ([#1749](https://github.com/rancher/terraform-provider-rancher2/issues/1749)) ([fe8c50c](https://github.com/rancher/terraform-provider-rancher2/commit/fe8c50c54251661e9d16a36b1671a50d058fa9e7))
* Add Oracle Kubernetes Engine back to terraform ([#1537](https://github.com/rancher/terraform-provider-rancher2/issues/1537)) ([00d15d7](https://github.com/rancher/terraform-provider-rancher2/commit/00d15d7ebd7b72c2c2c5d6a670226aaa2e3e3864))
* add permissions to update issues ([#1939](https://github.com/rancher/terraform-provider-rancher2/issues/1939)) ([8c19a61](https://github.com/rancher/terraform-provider-rancher2/commit/8c19a616d645bb19cb21fc81fef226af09cc1265))
* add recently added options to Oracle config ([328523d](https://github.com/rancher/terraform-provider-rancher2/commit/328523ddb75b773e5ae0e816310fc01ca55175e8))
* add support for data directories ([#1665](https://github.com/rancher/terraform-provider-rancher2/issues/1665)) ([#1684](https://github.com/rancher/terraform-provider-rancher2/issues/1684)) ([570d171](https://github.com/rancher/terraform-provider-rancher2/commit/570d171abc8b3cb54a1483e76adeea2b8f7873b3))
* add support for imported cluster config ([#1528](https://github.com/rancher/terraform-provider-rancher2/issues/1528)) ([e9e2669](https://github.com/rancher/terraform-provider-rancher2/commit/e9e2669896ffe76aebf5c0ce67816c2765351883))
* adjust path and remove module context ([#1706](https://github.com/rancher/terraform-provider-rancher2/issues/1706)) ([#1725](https://github.com/rancher/terraform-provider-rancher2/issues/1725)) ([df63d8e](https://github.com/rancher/terraform-provider-rancher2/commit/df63d8ee20ca8f4d190f736c561d689ee81dc9ba))
* allow specifying sha ([#1767](https://github.com/rancher/terraform-provider-rancher2/issues/1767)) ([3f3453b](https://github.com/rancher/terraform-provider-rancher2/commit/3f3453b9056523247240e1a35c3a25c21666be72))
* block use of deprecated Rke config ([535beb1](https://github.com/rancher/terraform-provider-rancher2/commit/535beb1387429d2df79d2f0c1236915c490e24eb))
* break out the config pkg for better testing ([#1902](https://github.com/rancher/terraform-provider-rancher2/issues/1902)) ([8a09d67](https://github.com/rancher/terraform-provider-rancher2/commit/8a09d6790b12f3417db17e8891b2234268ad1d77))
* build before testing ([#1906](https://github.com/rancher/terraform-provider-rancher2/issues/1906)) ([8dcb3e3](https://github.com/rancher/terraform-provider-rancher2/commit/8dcb3e35febe33d6c7245f84d3e338a17668baf0))
* change from "main" to "tracking" ([#1712](https://github.com/rancher/terraform-provider-rancher2/issues/1712)) ([#1728](https://github.com/rancher/terraform-provider-rancher2/issues/1728)) ([8d06a68](https://github.com/rancher/terraform-provider-rancher2/commit/8d06a6862715e734de983dcc8406d45332ec34f4))
* change IP for server URL in docs ([#1613](https://github.com/rancher/terraform-provider-rancher2/issues/1613)) ([338812e](https://github.com/rancher/terraform-provider-rancher2/commit/338812e8a4b170b86bb79c49970847e0a4f16687))
* change log from info to error ([#1604](https://github.com/rancher/terraform-provider-rancher2/issues/1604)) ([2e64ed4](https://github.com/rancher/terraform-provider-rancher2/commit/2e64ed4d4d69493b8924cc700b3c60c80a009b2d))
* check out repository ([#1754](https://github.com/rancher/terraform-provider-rancher2/issues/1754)) ([04e9bd7](https://github.com/rancher/terraform-provider-rancher2/commit/04e9bd7575f6c31c5622594228ace164df3aea2a))
* Clean up test suite ([#1587](https://github.com/rancher/terraform-provider-rancher2/issues/1587)) ([e8a3de0](https://github.com/rancher/terraform-provider-rancher2/commit/e8a3de02253548394437a8f3b5f146ff4dc4cc60))
* clear space on runner before release ([#1851](https://github.com/rancher/terraform-provider-rancher2/issues/1851)) ([c1925da](https://github.com/rancher/terraform-provider-rancher2/commit/c1925da9c2544d5bea261d2b924f8b4e8bcd5665))
* clear space on runners before release tests ([#1841](https://github.com/rancher/terraform-provider-rancher2/issues/1841)) ([d2cadf6](https://github.com/rancher/terraform-provider-rancher2/commit/d2cadf6afeaa20164a859e0ab8f4ff9318b94515))
* cluster proxy config ([#1737](https://github.com/rancher/terraform-provider-rancher2/issues/1737)) ([71c12ba](https://github.com/rancher/terraform-provider-rancher2/commit/71c12ba201fa8ca4c1145cfe49b19193effdc0b7))
* Correct genericoidc provider and tests ([#1931](https://github.com/rancher/terraform-provider-rancher2/issues/1931)) ([24d64d5](https://github.com/rancher/terraform-provider-rancher2/commit/24d64d5897e68b5e4a19fa01aaa738c32479f472))
* do not check out branch that we tag ([#1935](https://github.com/rancher/terraform-provider-rancher2/issues/1935)) ([58294d1](https://github.com/rancher/terraform-provider-rancher2/commit/58294d181b0a363b282d283d75062abcb1e2fdf5))
* do not check user state on creation ([#1734](https://github.com/rancher/terraform-provider-rancher2/issues/1734)) ([ab0fb60](https://github.com/rancher/terraform-provider-rancher2/commit/ab0fb600cf2ad509f9a141ed0033ffa354c5bcdb))
* do not use dot syntax in path ([#1705](https://github.com/rancher/terraform-provider-rancher2/issues/1705)) ([#1724](https://github.com/rancher/terraform-provider-rancher2/issues/1724)) ([4635e85](https://github.com/rancher/terraform-provider-rancher2/commit/4635e85516a534c2f86763e2e594dde42362b8c2))
* emit all response data for troubleshooting ([#1714](https://github.com/rancher/terraform-provider-rancher2/issues/1714)) ([#1729](https://github.com/rancher/terraform-provider-rancher2/issues/1729)) ([c55a983](https://github.com/rancher/terraform-provider-rancher2/commit/c55a9832b8ca80c3a62290414d807080948ddca2))
* enable manual rc release for older branches ([#1837](https://github.com/rancher/terraform-provider-rancher2/issues/1837)) ([ad5301f](https://github.com/rancher/terraform-provider-rancher2/commit/ad5301fb19ea28ae87dba8e06d900785d4f91bb2))
* ensure the correct workflow, users, and tags ([#1775](https://github.com/rancher/terraform-provider-rancher2/issues/1775)) ([ee9311f](https://github.com/rancher/terraform-provider-rancher2/commit/ee9311f4211fff29d0bd69d3f47eecaf13a61621))
* export GitHub scripts and validate in CI ([#1702](https://github.com/rancher/terraform-provider-rancher2/issues/1702)) ([#1722](https://github.com/rancher/terraform-provider-rancher2/issues/1722)) ([bacb1a1](https://github.com/rancher/terraform-provider-rancher2/commit/bacb1a1781f60d3a5291d34d6b87c942dcb34340))
* for manual releases skip the git validation ([#1791](https://github.com/rancher/terraform-provider-rancher2/issues/1791)) ([09b1ce0](https://github.com/rancher/terraform-provider-rancher2/commit/09b1ce06ac6c50da65b9ea3a425e05b61e0b2b72))
* force directory path ([#1710](https://github.com/rancher/terraform-provider-rancher2/issues/1710)) ([#1727](https://github.com/rancher/terraform-provider-rancher2/issues/1727)) ([b3a177c](https://github.com/rancher/terraform-provider-rancher2/commit/b3a177c61df93324222d673c29ddf7459993f11c))
* github script uses slashes for comments ([#1866](https://github.com/rancher/terraform-provider-rancher2/issues/1866)) ([96d0fcf](https://github.com/rancher/terraform-provider-rancher2/commit/96d0fcf52e13b5046a04a2ebf8bf12e4001a6bd9))
* goreleaser needs git state at the new tag ([#1771](https://github.com/rancher/terraform-provider-rancher2/issues/1771)) ([3a7a325](https://github.com/rancher/terraform-provider-rancher2/commit/3a7a325a5cc32a582c106548f1599eeb127c03fd))
* improve cleanup and skip in CI ([#1898](https://github.com/rancher/terraform-provider-rancher2/issues/1898)) ([098077d](https://github.com/rancher/terraform-provider-rancher2/commit/098077d1e06d102f44d16374809fcd1c6cd95099))
* improve registry documentation ([#1664](https://github.com/rancher/terraform-provider-rancher2/issues/1664)) ([#1679](https://github.com/rancher/terraform-provider-rancher2/issues/1679)) ([e7ef96a](https://github.com/rancher/terraform-provider-rancher2/commit/e7ef96a3aa6d60fea95c93396b704333255af997))
* improve test isolation logic ([#1882](https://github.com/rancher/terraform-provider-rancher2/issues/1882)) ([c59bf3c](https://github.com/rancher/terraform-provider-rancher2/commit/c59bf3ce3e22d9834ea622379375eebf65f64cfb))
* improve the documentation for cluster resources  ([#1821](https://github.com/rancher/terraform-provider-rancher2/issues/1821)) ([3f29e3a](https://github.com/rancher/terraform-provider-rancher2/commit/3f29e3a596492c22a6aa56c9d7e195a00b8dc663))
* make sure to only match the tag specified ([#1779](https://github.com/rancher/terraform-provider-rancher2/issues/1779)) ([6cc8710](https://github.com/rancher/terraform-provider-rancher2/commit/6cc8710d137bd2b4dc73f8fa7658f221cb77b6a5))
* manual full release workflow ([#1763](https://github.com/rancher/terraform-provider-rancher2/issues/1763)) ([a8ab38b](https://github.com/rancher/terraform-provider-rancher2/commit/a8ab38b2c8be770a2a88be9e6ec146e3754c11c0))
* missing quotes causing syntax error ([#1910](https://github.com/rancher/terraform-provider-rancher2/issues/1910)) ([d088af4](https://github.com/rancher/terraform-provider-rancher2/commit/d088af453768fe7b5efa0b74b7634569829b9d80))
* need to quote json list ([#1878](https://github.com/rancher/terraform-provider-rancher2/issues/1878)) ([967289f](https://github.com/rancher/terraform-provider-rancher2/commit/967289f3519f6359a6f1ee1eb0523efebcd3529d))
* notifications to issues for manual rc ([#1926](https://github.com/rancher/terraform-provider-rancher2/issues/1926)) ([55db4af](https://github.com/rancher/terraform-provider-rancher2/commit/55db4af8f93a0f7d8a96ee163a3b52368d2d300a))
* PSACT deletion hanging ([#1653](https://github.com/rancher/terraform-provider-rancher2/issues/1653)) ([#1672](https://github.com/rancher/terraform-provider-rancher2/issues/1672)) ([452fc1a](https://github.com/rancher/terraform-provider-rancher2/commit/452fc1aabe62c83990427e0190dbd51ab2a8565a))
* registry configuration for imported hosted clusters ([#1596](https://github.com/rancher/terraform-provider-rancher2/issues/1596)) ([826f513](https://github.com/rancher/terraform-provider-rancher2/commit/826f51349ea3af617600fb1046b341b04cc9015c))
* remove broken resources ([#1800](https://github.com/rancher/terraform-provider-rancher2/issues/1800)) ([783ccbd](https://github.com/rancher/terraform-provider-rancher2/commit/783ccbde487bcbddc6f779f459a4624e273da2c6))
* remove the default template ([#1783](https://github.com/rancher/terraform-provider-rancher2/issues/1783)) ([3fbfc6e](https://github.com/rancher/terraform-provider-rancher2/commit/3fbfc6ef28c64ea0e17b341a660930483c59f1c8))
* remove unused RKE version retrieval logic ([c9b2417](https://github.com/rancher/terraform-provider-rancher2/commit/c9b24176c89e8ed16c8f9180882bfc86ae82fad7))
* repository owner is in github context ([#1894](https://github.com/rancher/terraform-provider-rancher2/issues/1894)) ([433561d](https://github.com/rancher/terraform-provider-rancher2/commit/433561d7df3eb2697f812ee55e7d88103621e980))
* run tests in different jobs ([#1862](https://github.com/rancher/terraform-provider-rancher2/issues/1862)) ([4618c27](https://github.com/rancher/terraform-provider-rancher2/commit/4618c27906cba66ac97707c0d6580ce427df686c))
* Sanitize secrets ([#1568](https://github.com/rancher/terraform-provider-rancher2/issues/1568)) ([057d9cb](https://github.com/rancher/terraform-provider-rancher2/commit/057d9cbbc4dc907b8b8761fa04c00ed64c066119))
* script paths ([#1704](https://github.com/rancher/terraform-provider-rancher2/issues/1704)) ([#1723](https://github.com/rancher/terraform-provider-rancher2/issues/1723)) ([6a9c227](https://github.com/rancher/terraform-provider-rancher2/commit/6a9c227b9ffdc45f2d05c45bae027eb8649ec048))
* set bootstrap sha and release version ([#1758](https://github.com/rancher/terraform-provider-rancher2/issues/1758)) ([e42cccb](https://github.com/rancher/terraform-provider-rancher2/commit/e42cccb958dd8443f662841992806cd33a491669))
* short circuit the testing logic ([#1870](https://github.com/rancher/terraform-provider-rancher2/issues/1870)) ([04adaf2](https://github.com/rancher/terraform-provider-rancher2/commit/04adaf228f1d76f2d784251f0699b374916f20ff))
* specify signing key ([71a0c51](https://github.com/rancher/terraform-provider-rancher2/commit/71a0c51c9b892d51fa9ac1b8c60277f483a6b358))
* support exceptions in diff suppression ([#1826](https://github.com/rancher/terraform-provider-rancher2/issues/1826)) ([60f56e7](https://github.com/rancher/terraform-provider-rancher2/commit/60f56e7288aa96120e54f3cc3149cebee32bbf0f))
* the first job should never skip ([#1890](https://github.com/rancher/terraform-provider-rancher2/issues/1890)) ([376ed08](https://github.com/rancher/terraform-provider-rancher2/commit/376ed08bc2304bb78373a10b0b312fe70982ec77))
* typo in repository context ([#1886](https://github.com/rancher/terraform-provider-rancher2/issues/1886)) ([067b19e](https://github.com/rancher/terraform-provider-rancher2/commit/067b19e02337de502d70bfca4cb00fb593445339))
* typo in script name ([#1918](https://github.com/rancher/terraform-provider-rancher2/issues/1918)) ([c28e049](https://github.com/rancher/terraform-provider-rancher2/commit/c28e049c4564c850812ff14abccf9896813cd28e))
* typo in variable name ([#1914](https://github.com/rancher/terraform-provider-rancher2/issues/1914)) ([263bf4c](https://github.com/rancher/terraform-provider-rancher2/commit/263bf4ce678d6133170a85f9a32f2c23fa27d653))
* update checkout ([#1845](https://github.com/rancher/terraform-provider-rancher2/issues/1845)) ([5b7b307](https://github.com/rancher/terraform-provider-rancher2/commit/5b7b307acb2db1ee49a12bbfa76002e7c981ce6a))
* update label field to match latest version ([fc7015b](https://github.com/rancher/terraform-provider-rancher2/commit/fc7015bf0d46cb0c4690d6386dc67a2dd3880356))
* update pre-release to see new tag format ([#1623](https://github.com/rancher/terraform-provider-rancher2/issues/1623)) ([#1681](https://github.com/rancher/terraform-provider-rancher2/issues/1681)) ([2061872](https://github.com/rancher/terraform-provider-rancher2/commit/2061872a88b7240fdba187759d6bb2284cfa6dcd))
* update release config to get proper version ([#1958](https://github.com/rancher/terraform-provider-rancher2/issues/1958)) ([0382d8a](https://github.com/rancher/terraform-provider-rancher2/commit/0382d8a51ad77cf46e16672d4a2cf3dff3f3fcf1))
* update release manifest to force new version ([#1954](https://github.com/rancher/terraform-provider-rancher2/issues/1954)) ([16fa5a8](https://github.com/rancher/terraform-provider-rancher2/commit/16fa5a86c51557d8f9ef5e0fb740954572bf960b))
* update release to 13.1.0 ([#1946](https://github.com/rancher/terraform-provider-rancher2/issues/1946)) ([87fe193](https://github.com/rancher/terraform-provider-rancher2/commit/87fe193a07fc8cfd254285fc1d09c0d5721856b2))
* update release to 13.1.1 ([#1950](https://github.com/rancher/terraform-provider-rancher2/issues/1950)) ([43564c9](https://github.com/rancher/terraform-provider-rancher2/commit/43564c9632e881fe22eced139da1f49d6b45a3ec))
* use back tick instead of single quotes ([#1855](https://github.com/rancher/terraform-provider-rancher2/issues/1855)) ([414ede7](https://github.com/rancher/terraform-provider-rancher2/commit/414ede70c5a76d23c823d9b49e202ba6f0b72529))
* use environment to pass data to the script ([#1922](https://github.com/rancher/terraform-provider-rancher2/issues/1922)) ([5d58a8a](https://github.com/rancher/terraform-provider-rancher2/commit/5d58a8a018e68571e7500ed9191d83f1cdfa5305))
* use modern import logic ([#1708](https://github.com/rancher/terraform-provider-rancher2/issues/1708)) ([#1726](https://github.com/rancher/terraform-provider-rancher2/issues/1726)) ([d435b19](https://github.com/rancher/terraform-provider-rancher2/commit/d435b198653ba558f9e7f6736a5c3e1d34f5db8f))
* use proper outputs for validating test ([#1874](https://github.com/rancher/terraform-provider-rancher2/issues/1874)) ([b7289ce](https://github.com/rancher/terraform-provider-rancher2/commit/b7289ce056d1f48e47aba5b64711c5d5e3c68254))
* use separate goreleaser files for manual ([#1787](https://github.com/rancher/terraform-provider-rancher2/issues/1787)) ([1886bc3](https://github.com/rancher/terraform-provider-rancher2/commit/1886bc3149d45f1377bfd8d9e65674fe89b2ab6a))
* validate connection ([#1585](https://github.com/rancher/terraform-provider-rancher2/issues/1585)) ([1fad9c4](https://github.com/rancher/terraform-provider-rancher2/commit/1fad9c4f4069fe5c56d6a33be10e417fe4dee5c4))

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
