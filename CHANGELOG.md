# 📦 open-sauced/pizza-cli changelog

[![conventional commits](https://img.shields.io/badge/conventional%20commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![semantic versioning](https://img.shields.io/badge/semantic%20versioning-2.0.0-green.svg)](https://semver.org)

> All notable changes to this project will be documented in this file

## [1.4.0-beta.4](https://github.com/open-sauced/pizza-cli/compare/v1.4.0-beta.3...v1.4.0-beta.4) (2024-09-09)


### 🐛 Bug Fixes

* docs generation runs outside of build matrix now ([#165](https://github.com/open-sauced/pizza-cli/issues/165)) ([1e42988](https://github.com/open-sauced/pizza-cli/commit/1e42988c06fcab6694d4fca9670c59796352e7a5))

## [1.4.0-beta.3](https://github.com/open-sauced/pizza-cli/compare/v1.4.0-beta.2...v1.4.0-beta.3) (2024-09-09)


### 🐛 Bug Fixes

* now --tty-disable is set so the action can complete instead of hanging ([#164](https://github.com/open-sauced/pizza-cli/issues/164)) ([a970a73](https://github.com/open-sauced/pizza-cli/commit/a970a73f494f34464a4c8b6ba993d38ecb4e2ec4))

## [1.4.0-beta.2](https://github.com/open-sauced/pizza-cli/compare/v1.4.0-beta.1...v1.4.0-beta.2) (2024-09-09)


### 🐛 Bug Fixes

* fixed docs generation in release workflow ([#162](https://github.com/open-sauced/pizza-cli/issues/162)) ([5341e16](https://github.com/open-sauced/pizza-cli/commit/5341e16daaeeecdc664895d165246a82623accbe))

## [1.4.0-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.3.1-beta.2...v1.4.0-beta.1) (2024-09-09)


### 🍕 Features

* Posthog events bootstrapping ([#160](https://github.com/open-sauced/pizza-cli/issues/160)) ([847426b](https://github.com/open-sauced/pizza-cli/commit/847426bcb202e8846287461fb0e3735d04f4c82e))

## [1.3.1-beta.2](https://github.com/open-sauced/pizza-cli/compare/v1.3.1-beta.1...v1.3.1-beta.2) (2024-09-06)


### 🐛 Bug Fixes

* use the local directory and home directory as fallback for .sauced.yaml ([#158](https://github.com/open-sauced/pizza-cli/issues/158)) ([af2f361](https://github.com/open-sauced/pizza-cli/commit/af2f3612e26634455602d1840714c5bf15e1e40a))

## [1.3.1-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.3.0...v1.3.1-beta.1) (2024-09-06)


### 🐛 Bug Fixes

* skip interactive steps in generate codeowners with --tty-disable flag ([#159](https://github.com/open-sauced/pizza-cli/issues/159)) ([49f1fd3](https://github.com/open-sauced/pizza-cli/commit/49f1fd3fc4df24b95724feb1918dc80276cd017e))

## [1.3.0](https://github.com/open-sauced/pizza-cli/compare/v1.2.1...v1.3.0) (2024-09-06)


### 🍕 Features

* Create a contributor list after generating codeowners ([#141](https://github.com/open-sauced/pizza-cli/issues/141)) ([72c5d58](https://github.com/open-sauced/pizza-cli/commit/72c5d588fcd4fb04f6d39d756a6a26a47d25a4e4))
* now the documentation for the pizza-cli can be generated via pizza docs ([#143](https://github.com/open-sauced/pizza-cli/issues/143)) ([3f5d27e](https://github.com/open-sauced/pizza-cli/commit/3f5d27e2c52c894a266828e70e7475069e74e8e9))
* Refactors API client into hand rolled sdk in api/ directory ([#111](https://github.com/open-sauced/pizza-cli/issues/111)) ([e16e889](https://github.com/open-sauced/pizza-cli/commit/e16e8899a4ef69641dc614887d065dc8b70adb35))
* support fallback attributions when generating codeowners file ([#145](https://github.com/open-sauced/pizza-cli/issues/145)) ([35af4da](https://github.com/open-sauced/pizza-cli/commit/35af4dafc4ed088ba1396ff28e1536723c914a2b))
* update `CODEOWNERS` copy with command ([#130](https://github.com/open-sauced/pizza-cli/issues/130)) ([a477959](https://github.com/open-sauced/pizza-cli/commit/a477959020cfcbb3dc4707efb1700e17e05e3981))


### 🐛 Bug Fixes

* Corrects invalid gosec lint error ([#151](https://github.com/open-sauced/pizza-cli/issues/151)) ([f76527f](https://github.com/open-sauced/pizza-cli/commit/f76527f0c61c5720f684416f391fe1395774e1fb))
* Exhume Posthog functionality ([#147](https://github.com/open-sauced/pizza-cli/issues/147)) ([de091ca](https://github.com/open-sauced/pizza-cli/commit/de091cac7df585eadcfae64d6f851cfc178c74a2))
* now fallback .sauced.yaml contents get read ([#135](https://github.com/open-sauced/pizza-cli/issues/135)) ([fd658e5](https://github.com/open-sauced/pizza-cli/commit/fd658e5e09051cdf007c3605aa880d68db835afb))
* NPM cache now looks at package-lock file ([#136](https://github.com/open-sauced/pizza-cli/issues/136)) ([cd4b8da](https://github.com/open-sauced/pizza-cli/commit/cd4b8da75e1a0c0aa3d7e6f76d6b560a4dea941f))

## [1.3.0-beta.9](https://github.com/open-sauced/pizza-cli/compare/v1.3.0-beta.8...v1.3.0-beta.9) (2024-09-06)


### 🐛 Bug Fixes

* Corrects invalid gosec lint error ([#151](https://github.com/open-sauced/pizza-cli/issues/151)) ([f76527f](https://github.com/open-sauced/pizza-cli/commit/f76527f0c61c5720f684416f391fe1395774e1fb))

## [1.3.0-beta.8](https://github.com/open-sauced/pizza-cli/compare/v1.3.0-beta.7...v1.3.0-beta.8) (2024-09-06)


### 🍕 Features

* now the documentation for the pizza-cli can be generated via pizza docs ([#143](https://github.com/open-sauced/pizza-cli/issues/143)) ([3f5d27e](https://github.com/open-sauced/pizza-cli/commit/3f5d27e2c52c894a266828e70e7475069e74e8e9))

## [1.3.0-beta.7](https://github.com/open-sauced/pizza-cli/compare/v1.3.0-beta.6...v1.3.0-beta.7) (2024-09-06)


### 🐛 Bug Fixes

* Exhume Posthog functionality ([#147](https://github.com/open-sauced/pizza-cli/issues/147)) ([de091ca](https://github.com/open-sauced/pizza-cli/commit/de091cac7df585eadcfae64d6f851cfc178c74a2))

## [1.3.0-beta.6](https://github.com/open-sauced/pizza-cli/compare/v1.3.0-beta.5...v1.3.0-beta.6) (2024-09-05)


### 🍕 Features

* support fallback attributions when generating codeowners file ([#145](https://github.com/open-sauced/pizza-cli/issues/145)) ([35af4da](https://github.com/open-sauced/pizza-cli/commit/35af4dafc4ed088ba1396ff28e1536723c914a2b))

## [1.3.0-beta.5](https://github.com/open-sauced/pizza-cli/compare/v1.3.0-beta.4...v1.3.0-beta.5) (2024-09-05)


### 🍕 Features

* Create a contributor list after generating codeowners ([#141](https://github.com/open-sauced/pizza-cli/issues/141)) ([72c5d58](https://github.com/open-sauced/pizza-cli/commit/72c5d588fcd4fb04f6d39d756a6a26a47d25a4e4))

## [1.3.0-beta.4](https://github.com/open-sauced/pizza-cli/compare/v1.3.0-beta.3...v1.3.0-beta.4) (2024-09-05)


### 🐛 Bug Fixes

* now fallback .sauced.yaml contents get read ([#135](https://github.com/open-sauced/pizza-cli/issues/135)) ([fd658e5](https://github.com/open-sauced/pizza-cli/commit/fd658e5e09051cdf007c3605aa880d68db835afb))

## [1.3.0-beta.3](https://github.com/open-sauced/pizza-cli/compare/v1.3.0-beta.2...v1.3.0-beta.3) (2024-09-04)


### 🐛 Bug Fixes

* NPM cache now looks at package-lock file ([#136](https://github.com/open-sauced/pizza-cli/issues/136)) ([cd4b8da](https://github.com/open-sauced/pizza-cli/commit/cd4b8da75e1a0c0aa3d7e6f76d6b560a4dea941f))

## [1.3.0-beta.2](https://github.com/open-sauced/pizza-cli/compare/v1.3.0-beta.1...v1.3.0-beta.2) (2024-09-04)


### 🍕 Features

* update `CODEOWNERS` copy with command ([#130](https://github.com/open-sauced/pizza-cli/issues/130)) ([a477959](https://github.com/open-sauced/pizza-cli/commit/a477959020cfcbb3dc4707efb1700e17e05e3981))

## [1.3.0-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.2.1...v1.3.0-beta.1) (2024-09-04)


### 🍕 Features

* Refactors API client into hand rolled sdk in api/ directory ([#111](https://github.com/open-sauced/pizza-cli/issues/111)) ([e16e889](https://github.com/open-sauced/pizza-cli/commit/e16e8899a4ef69641dc614887d065dc8b70adb35))

## [1.2.1](https://github.com/open-sauced/pizza-cli/compare/v1.2.0...v1.2.1) (2024-08-30)


### 🐛 Bug Fixes

* Root command persistent flags are marked hidden correctly ([#126](https://github.com/open-sauced/pizza-cli/issues/126)) ([727a82e](https://github.com/open-sauced/pizza-cli/commit/727a82e6488699ae854df6d7dc1ac0778ef03542))

## [1.2.1-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.2.0...v1.2.1-beta.1) (2024-08-30)


### 🐛 Bug Fixes

* Root command persistent flags are marked hidden correctly ([#126](https://github.com/open-sauced/pizza-cli/issues/126)) ([727a82e](https://github.com/open-sauced/pizza-cli/commit/727a82e6488699ae854df6d7dc1ac0778ef03542))

## [1.2.0](https://github.com/open-sauced/pizza-cli/compare/v1.1.1...v1.2.0) (2024-08-30)


### 📝 Documentation

* updated comment for LoadConfig ([ab5206b](https://github.com/open-sauced/pizza-cli/commit/ab5206b9a76ca35e2cf18c0fa68c4958a7d37ca6))


### ✅ Tests

* added setup and teardown ([aba6310](https://github.com/open-sauced/pizza-cli/commit/aba631095495fc70b2ba520a1e59e5d7c6f93a13))
* added tests for LoadConfig ([2a5f85d](https://github.com/open-sauced/pizza-cli/commit/2a5f85de3bbe3a3d2812fcfe1ba4640ad7b55827))
* made tests parallel ([2b3d8ca](https://github.com/open-sauced/pizza-cli/commit/2b3d8ca28e7a92f81953e16a982acda5a36a4e6d))


### 🍕 Features

* added .sauced.yaml and updated CODEOWNERS file ([#109](https://github.com/open-sauced/pizza-cli/issues/109)) ([dfc56cb](https://github.com/open-sauced/pizza-cli/commit/dfc56cbafd0f061bc0742e3f8a1d8be93ed7bfda))
* added built at to version command ([#94](https://github.com/open-sauced/pizza-cli/issues/94)) ([9960fc0](https://github.com/open-sauced/pizza-cli/commit/9960fc0733e3f6c22692b0bc89ff00c674a97274))
* Codeowners generation ([#95](https://github.com/open-sauced/pizza-cli/issues/95)) ([79cf8a2](https://github.com/open-sauced/pizza-cli/commit/79cf8a2b47c701505bd889df569c592bfed49dbd))
* now generate codeowners checks in user root but also repository root folder for .sauced.yaml ([a0298b1](https://github.com/open-sauced/pizza-cli/commit/a0298b1b8bd0348918d579785cc859cefa594ada))
* pizza login success page style refresh ([#112](https://github.com/open-sauced/pizza-cli/issues/112)) ([9357dac](https://github.com/open-sauced/pizza-cli/commit/9357dac1bf07cc3459cae100e78ecdf451747544))
* Refactors Auth code into api/ directory ([#105](https://github.com/open-sauced/pizza-cli/issues/105)) ([d851499](https://github.com/open-sauced/pizza-cli/commit/d851499690b0038d10a17ea67019021ffe7c70f2))
* Skip semantic-release docker build in favor of buildx building ([d782974](https://github.com/open-sauced/pizza-cli/commit/d782974021739f2e23d95d1f9f35ff01b24b628b))
* Trim down CLI - remove unused, defunct commands ([#93](https://github.com/open-sauced/pizza-cli/issues/93)) ([7ddd4b9](https://github.com/open-sauced/pizza-cli/commit/7ddd4b971eb085b995db8205f7fc718701ec4db4))
* Upgrade Go module to use Go 1.22 ([#96](https://github.com/open-sauced/pizza-cli/issues/96)) ([690b6e9](https://github.com/open-sauced/pizza-cli/commit/690b6e92549b4478e0d6a8cf1814052158c851cd))
* Use justfile vs. makefile ([#84](https://github.com/open-sauced/pizza-cli/issues/84)) ([8f38eaf](https://github.com/open-sauced/pizza-cli/commit/8f38eaf4f24947a4035b5ac764899acc55748d38))


### 🐛 Bug Fixes

* escape non-standard characters in the filename path ([#106](https://github.com/open-sauced/pizza-cli/issues/106)) ([418951f](https://github.com/open-sauced/pizza-cli/commit/418951f2629412c0855161a82eae1fd87502091b))
* mark endpoint and beta flags as hidden ([#113](https://github.com/open-sauced/pizza-cli/issues/113)) ([6aa250f](https://github.com/open-sauced/pizza-cli/commit/6aa250f45bde6486189270da2c70589f475ca39c))
* move output flag to insights command ([#115](https://github.com/open-sauced/pizza-cli/issues/115)) ([be7f8cd](https://github.com/open-sauced/pizza-cli/commit/be7f8cdda34552c5f845e995c84a3a1aa4be01d2))
* update supabase keys for auth ([#80](https://github.com/open-sauced/pizza-cli/issues/80)) ([247c431](https://github.com/open-sauced/pizza-cli/commit/247c431401744833e36ebea797c49062f0e35910))
* use repository fullname to fetch contributors ([#77](https://github.com/open-sauced/pizza-cli/issues/77)) ([5326875](https://github.com/open-sauced/pizza-cli/commit/53268758056a25b9135c011425b5854657752885))

## [1.2.0-beta.12](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.11...v1.2.0-beta.12) (2024-08-30)


### 🐛 Bug Fixes

* move output flag to insights command ([#115](https://github.com/open-sauced/pizza-cli/issues/115)) ([be7f8cd](https://github.com/open-sauced/pizza-cli/commit/be7f8cdda34552c5f845e995c84a3a1aa4be01d2))

## [1.2.0-beta.11](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.10...v1.2.0-beta.11) (2024-08-29)


### 🐛 Bug Fixes

* mark endpoint and beta flags as hidden ([#113](https://github.com/open-sauced/pizza-cli/issues/113)) ([6aa250f](https://github.com/open-sauced/pizza-cli/commit/6aa250f45bde6486189270da2c70589f475ca39c))

## [1.2.0-beta.10](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.9...v1.2.0-beta.10) (2024-08-29)


### 🍕 Features

* pizza login success page style refresh ([#112](https://github.com/open-sauced/pizza-cli/issues/112)) ([9357dac](https://github.com/open-sauced/pizza-cli/commit/9357dac1bf07cc3459cae100e78ecdf451747544))

## [1.2.0-beta.9](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.8...v1.2.0-beta.9) (2024-08-29)


### 🍕 Features

* now generate codeowners checks in user root but also repository root folder for .sauced.yaml ([a0298b1](https://github.com/open-sauced/pizza-cli/commit/a0298b1b8bd0348918d579785cc859cefa594ada))


### 📝 Documentation

* updated comment for LoadConfig ([ab5206b](https://github.com/open-sauced/pizza-cli/commit/ab5206b9a76ca35e2cf18c0fa68c4958a7d37ca6))


### ✅ Tests

* added setup and teardown ([aba6310](https://github.com/open-sauced/pizza-cli/commit/aba631095495fc70b2ba520a1e59e5d7c6f93a13))
* added tests for LoadConfig ([2a5f85d](https://github.com/open-sauced/pizza-cli/commit/2a5f85de3bbe3a3d2812fcfe1ba4640ad7b55827))
* made tests parallel ([2b3d8ca](https://github.com/open-sauced/pizza-cli/commit/2b3d8ca28e7a92f81953e16a982acda5a36a4e6d))

## [1.2.0-beta.8](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.7...v1.2.0-beta.8) (2024-08-28)


### 🍕 Features

* added .sauced.yaml and updated CODEOWNERS file ([#109](https://github.com/open-sauced/pizza-cli/issues/109)) ([dfc56cb](https://github.com/open-sauced/pizza-cli/commit/dfc56cbafd0f061bc0742e3f8a1d8be93ed7bfda))

## [1.2.0-beta.7](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.6...v1.2.0-beta.7) (2024-08-28)


### 🍕 Features

* Refactors Auth code into api/ directory ([#105](https://github.com/open-sauced/pizza-cli/issues/105)) ([d851499](https://github.com/open-sauced/pizza-cli/commit/d851499690b0038d10a17ea67019021ffe7c70f2))

## [1.2.0-beta.6](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.5...v1.2.0-beta.6) (2024-08-28)


### 🐛 Bug Fixes

* escape non-standard characters in the filename path ([#106](https://github.com/open-sauced/pizza-cli/issues/106)) ([418951f](https://github.com/open-sauced/pizza-cli/commit/418951f2629412c0855161a82eae1fd87502091b))

## [1.2.0-beta.5](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.4...v1.2.0-beta.5) (2024-08-27)


### 🍕 Features

* added built at to version command ([#94](https://github.com/open-sauced/pizza-cli/issues/94)) ([9960fc0](https://github.com/open-sauced/pizza-cli/commit/9960fc0733e3f6c22692b0bc89ff00c674a97274))

## [1.2.0-beta.4](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.3...v1.2.0-beta.4) (2024-08-27)


### 🍕 Features

* Upgrade Go module to use Go 1.22 ([#96](https://github.com/open-sauced/pizza-cli/issues/96)) ([690b6e9](https://github.com/open-sauced/pizza-cli/commit/690b6e92549b4478e0d6a8cf1814052158c851cd))

## [1.2.0-beta.3](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.2...v1.2.0-beta.3) (2024-08-27)


### 🍕 Features

* Codeowners generation ([#95](https://github.com/open-sauced/pizza-cli/issues/95)) ([79cf8a2](https://github.com/open-sauced/pizza-cli/commit/79cf8a2b47c701505bd889df569c592bfed49dbd))

## [1.2.0-beta.2](https://github.com/open-sauced/pizza-cli/compare/v1.2.0-beta.1...v1.2.0-beta.2) (2024-08-26)


### 🍕 Features

* Trim down CLI - remove unused, defunct commands ([#93](https://github.com/open-sauced/pizza-cli/issues/93)) ([7ddd4b9](https://github.com/open-sauced/pizza-cli/commit/7ddd4b971eb085b995db8205f7fc718701ec4db4))

## [1.2.0-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.1.1-beta.5...v1.2.0-beta.1) (2024-08-26)


### 🍕 Features

* Skip semantic-release docker build in favor of buildx building ([d782974](https://github.com/open-sauced/pizza-cli/commit/d782974021739f2e23d95d1f9f35ff01b24b628b))
* Use justfile vs. makefile ([#84](https://github.com/open-sauced/pizza-cli/issues/84)) ([8f38eaf](https://github.com/open-sauced/pizza-cli/commit/8f38eaf4f24947a4035b5ac764899acc55748d38))

## [1.2.0-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.1.1-beta.5...v1.2.0-beta.1) (2024-08-26)


### 🍕 Features

* Use justfile vs. makefile ([#84](https://github.com/open-sauced/pizza-cli/issues/84)) ([8f38eaf](https://github.com/open-sauced/pizza-cli/commit/8f38eaf4f24947a4035b5ac764899acc55748d38))

## [1.1.1-beta.5](https://github.com/open-sauced/pizza-cli/compare/v1.1.1-beta.4...v1.1.1-beta.5) (2024-03-15)


### 🐛 Bug Fixes

* update supabase keys for auth ([#80](https://github.com/open-sauced/pizza-cli/issues/80)) ([247c431](https://github.com/open-sauced/pizza-cli/commit/247c431401744833e36ebea797c49062f0e35910))

## [1.1.1-beta.4](https://github.com/open-sauced/pizza-cli/compare/v1.1.1-beta.3...v1.1.1-beta.4) (2024-03-04)


### 🐛 Bug Fixes

* use repository fullname to fetch contributors ([#77](https://github.com/open-sauced/pizza-cli/issues/77)) ([5326875](https://github.com/open-sauced/pizza-cli/commit/53268758056a25b9135c011425b5854657752885))
* avoid requiring CGO for now ([#71](https://github.com/open-sauced/pizza-cli/issues/71)) ([f6d2f1d](https://github.com/open-sauced/pizza-cli/commit/f6d2f1d11bda7760edf585279099f3a874661973))
* Force publish of package to build go binaries during release ([#69](https://github.com/open-sauced/pizza-cli/issues/69)) ([02177d5](https://github.com/open-sauced/pizza-cli/commit/02177d5c81c330385f4f73c5d5f2df045c96757e))
* Upgrade to v2 API ([#73](https://github.com/open-sauced/pizza-cli/issues/73)) ([40b468b](https://github.com/open-sauced/pizza-cli/commit/40b468be69bdffb1fa7170861abf98601acb6c68))


## [1.1.1-beta.3](https://github.com/open-sauced/pizza-cli/compare/v1.1.1-beta.2...v1.1.1-beta.3) (2024-02-05)


### 🐛 Bug Fixes

* Upgrade to v2 API ([#73](https://github.com/open-sauced/pizza-cli/issues/73)) ([40b468b](https://github.com/open-sauced/pizza-cli/commit/40b468be69bdffb1fa7170861abf98601acb6c68))

## [1.1.1-beta.2](https://github.com/open-sauced/pizza-cli/compare/v1.1.1-beta.1...v1.1.1-beta.2) (2024-02-03)


### 🐛 Bug Fixes

* avoid requiring CGO for now ([#71](https://github.com/open-sauced/pizza-cli/issues/71)) ([f6d2f1d](https://github.com/open-sauced/pizza-cli/commit/f6d2f1d11bda7760edf585279099f3a874661973))

## [1.1.1-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.1.0...v1.1.1-beta.1) (2023-11-29)


### 🐛 Bug Fixes

* Force publish of package to build go binaries during release ([#69](https://github.com/open-sauced/pizza-cli/issues/69)) ([02177d5](https://github.com/open-sauced/pizza-cli/commit/02177d5c81c330385f4f73c5d5f2df045c96757e))

## [1.1.0](https://github.com/open-sauced/pizza-cli/compare/v1.0.1...v1.1.0) (2023-10-26)


### 🐛 Bug Fixes

* Hotfix for broken release workflow file ([e115fe9](https://github.com/open-sauced/pizza-cli/commit/e115fe91ded711ad075e2847b2c73f3f065b50d1))
* Release pipeline missing some env vars ([21e43e3](https://github.com/open-sauced/pizza-cli/commit/21e43e365d490eb0a9182895aaa5e45a2c00a025))


### 🍕 Features

* Add "pizza insights user-contributions" command ([248f90a](https://github.com/open-sauced/pizza-cli/commit/248f90a270540c7dd992b9d69dff7da1eb2ceddd))
* add csv support for contributors insights ([51c1897](https://github.com/open-sauced/pizza-cli/commit/51c18978cb9012fa72642fcb77c9f33bf0a88076))
* Add filtering for usernames on user-contributions ([d10c1ef](https://github.com/open-sauced/pizza-cli/commit/d10c1ef0bc69d3acd2e9f3ebdab7eee813debe19))

## [1.1.0-beta.3](https://github.com/open-sauced/pizza-cli/compare/v1.1.0-beta.2...v1.1.0-beta.3) (2023-10-26)


### 🍕 Features

* Add filtering for usernames on user-contributions ([d10c1ef](https://github.com/open-sauced/pizza-cli/commit/d10c1ef0bc69d3acd2e9f3ebdab7eee813debe19))

## [1.1.0-beta.2](https://github.com/open-sauced/pizza-cli/compare/v1.1.0-beta.1...v1.1.0-beta.2) (2023-10-26)


### 🍕 Features

* Add "pizza insights user-contributions" command ([248f90a](https://github.com/open-sauced/pizza-cli/commit/248f90a270540c7dd992b9d69dff7da1eb2ceddd))

## [1.1.0-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.0.1-beta.1...v1.1.0-beta.1) (2023-10-18)


### 🍕 Features

* add csv support for contributors insights ([51c1897](https://github.com/open-sauced/pizza-cli/commit/51c18978cb9012fa72642fcb77c9f33bf0a88076))

## [1.0.1](https://github.com/open-sauced/pizza-cli/compare/v1.0.0...v1.0.1) (2023-10-11)


### 🐛 Bug Fixes

* Hotfix for broken release workflow file ([e115fe9](https://github.com/open-sauced/pizza-cli/commit/e115fe91ded711ad075e2847b2c73f3f065b50d1))
* Release pipeline missing some env vars ([21e43e3](https://github.com/open-sauced/pizza-cli/commit/21e43e365d490eb0a9182895aaa5e45a2c00a025))
* Roll forward fix for semantic release ([#51](https://github.com/open-sauced/pizza-cli/issues/51)) ([3a6bb27](https://github.com/open-sauced/pizza-cli/commit/3a6bb27da209caccc4ab092e202516442b1cc621))

## [1.0.1-beta.1](https://github.com/open-sauced/pizza-cli/compare/v1.0.0...v1.0.1-beta.1) (2023-10-11)


### 🐛 Bug Fixes

* Roll forward fix for semantic release ([#51](https://github.com/open-sauced/pizza-cli/issues/51)) ([3a6bb27](https://github.com/open-sauced/pizza-cli/commit/3a6bb27da209caccc4ab092e202516442b1cc621))

## 1.0.0 (2023-10-11)


### 🤖 Build System

* sematic bin release, npm ([7b4607e](https://github.com/open-sauced/pizza-cli/commit/7b4607e9a4aa5eba0b5f163c586520c1022494ee))


### 🔁 Continuous Integration

* Update @open-sauced/release@2.2.1 and compliance.yaml ([#33](https://github.com/open-sauced/pizza-cli/issues/33)) ([146b6b7](https://github.com/open-sauced/pizza-cli/commit/146b6b7485a0f33090a4ccefd23624f9aa0df085))


### 🐛 Bug Fixes

* Uses correct generated token when checking out cli repo in release ([#44](https://github.com/open-sauced/pizza-cli/issues/44)) ([1e0c9f1](https://github.com/open-sauced/pizza-cli/commit/1e0c9f1ef3c9d0d9bd7f590f6bec021707f4c833))


### 🍕 Features

* Add install instructions and script for pizza CLI ([#26](https://github.com/open-sauced/pizza-cli/issues/26)) ([421a429](https://github.com/open-sauced/pizza-cli/commit/421a429ed99cca957365106485da97e085b0f173))
* Add posthog telemetry integration ([#37](https://github.com/open-sauced/pizza-cli/issues/37)) ([9829f49](https://github.com/open-sauced/pizza-cli/commit/9829f499dad0651ec97d0969e040d2acc30714e0))
* cli auth ([#21](https://github.com/open-sauced/pizza-cli/issues/21)) ([34728fb](https://github.com/open-sauced/pizza-cli/commit/34728fb62d01b746ffc8ede3c97a090b32b0b9f9))
* GitHub action to build and upload Go artifacts after release created ([#22](https://github.com/open-sauced/pizza-cli/issues/22)) ([ad187a9](https://github.com/open-sauced/pizza-cli/commit/ad187a9f3229e41785a09130132a799378c04528))
* Http Client for accessing OpenSauced API client ([#23](https://github.com/open-sauced/pizza-cli/issues/23)) ([ec2b357](https://github.com/open-sauced/pizza-cli/commit/ec2b35789a2864d38bf63e0ec1a3b68393a34e9b))
* Leverage the GITHUB_APP_TOKEN for releases ([#32](https://github.com/open-sauced/pizza-cli/issues/32)) ([e0a25e0](https://github.com/open-sauced/pizza-cli/commit/e0a25e003e89a7a5173ecaae12366922365243c9))
* npm i -g pizza ([73291d1](https://github.com/open-sauced/pizza-cli/commit/73291d13d632b709f2583d834aefe6ad758de8d7))
* Pizza show ([#24](https://github.com/open-sauced/pizza-cli/issues/24)) ([72f21ce](https://github.com/open-sauced/pizza-cli/commit/72f21ce260ec73c3ea0d7e97ed1411a86bb1d753))
* provide repository contributors insights ([#30](https://github.com/open-sauced/pizza-cli/issues/30)) ([d16091f](https://github.com/open-sauced/pizza-cli/commit/d16091ff4ee2ad74e025779b27321897d2c8a49c))
* provide repository insights ([#38](https://github.com/open-sauced/pizza-cli/issues/38)) ([dc148d6](https://github.com/open-sauced/pizza-cli/commit/dc148d6fe17b9aa96ad6951aefd0a7fd7cf0e160))
* repo-query support ([199cfd7](https://github.com/open-sauced/pizza-cli/commit/199cfd7b04e1e1683cce5abc08c57bbef01644f6))
* update bin name release.yaml ([6b21cb8](https://github.com/open-sauced/pizza-cli/commit/6b21cb84f88f75467ce6f270e136dfca5e462d23))
* Version command for CLI based on release builds ([#36](https://github.com/open-sauced/pizza-cli/issues/36)) ([9f3eedc](https://github.com/open-sauced/pizza-cli/commit/9f3eedcf7dac1d72f91ae40a5a09df3ee341a99c))

## [1.0.0-beta.7](https://github.com/open-sauced/pizza-cli/compare/v1.0.0-beta.6...v1.0.0-beta.7) (2023-10-02)


### 🍕 Features

* provide repository insights ([#38](https://github.com/open-sauced/pizza-cli/issues/38)) ([dc148d6](https://github.com/open-sauced/pizza-cli/commit/dc148d6fe17b9aa96ad6951aefd0a7fd7cf0e160))

## [1.0.0-beta.6](https://github.com/open-sauced/pizza-cli/compare/v1.0.0-beta.5...v1.0.0-beta.6) (2023-09-27)


### 🐛 Bug Fixes

* Uses correct generated token when checking out cli repo in release ([#44](https://github.com/open-sauced/pizza-cli/issues/44)) ([1e0c9f1](https://github.com/open-sauced/pizza-cli/commit/1e0c9f1ef3c9d0d9bd7f590f6bec021707f4c833))

## [1.0.0-beta.5](https://github.com/open-sauced/pizza-cli/compare/v1.0.0-beta.4...v1.0.0-beta.5) (2023-09-27)


### 🍕 Features

* Pizza show ([#24](https://github.com/open-sauced/pizza-cli/issues/24)) ([72f21ce](https://github.com/open-sauced/pizza-cli/commit/72f21ce260ec73c3ea0d7e97ed1411a86bb1d753))

## [1.0.0-beta.4](https://github.com/open-sauced/pizza-cli/compare/v1.0.0-beta.3...v1.0.0-beta.4) (2023-09-06)


### 🍕 Features

* provide repository contributors insights ([#30](https://github.com/open-sauced/pizza-cli/issues/30)) ([d16091f](https://github.com/open-sauced/pizza-cli/commit/d16091ff4ee2ad74e025779b27321897d2c8a49c))

## [1.0.0-beta.3](https://github.com/open-sauced/pizza-cli/compare/v1.0.0-beta.2...v1.0.0-beta.3) (2023-08-31)


### 🍕 Features

* Version command for CLI based on release builds ([#36](https://github.com/open-sauced/pizza-cli/issues/36)) ([9f3eedc](https://github.com/open-sauced/pizza-cli/commit/9f3eedcf7dac1d72f91ae40a5a09df3ee341a99c))

## [1.0.0-beta.2](https://github.com/open-sauced/pizza-cli/compare/v1.0.0-beta.1...v1.0.0-beta.2) (2023-08-29)


### 🍕 Features

* Add posthog telemetry integration ([#37](https://github.com/open-sauced/pizza-cli/issues/37)) ([9829f49](https://github.com/open-sauced/pizza-cli/commit/9829f499dad0651ec97d0969e040d2acc30714e0))

## 1.0.0-beta.1 (2023-08-22)


### 🤖 Build System

* sematic bin release, npm ([7b4607e](https://github.com/open-sauced/pizza-cli/commit/7b4607e9a4aa5eba0b5f163c586520c1022494ee))


### 🍕 Features

* Add install instructions and script for pizza CLI ([#26](https://github.com/open-sauced/pizza-cli/issues/26)) ([421a429](https://github.com/open-sauced/pizza-cli/commit/421a429ed99cca957365106485da97e085b0f173))
* cli auth ([#21](https://github.com/open-sauced/pizza-cli/issues/21)) ([34728fb](https://github.com/open-sauced/pizza-cli/commit/34728fb62d01b746ffc8ede3c97a090b32b0b9f9))
* GitHub action to build and upload Go artifacts after release created ([#22](https://github.com/open-sauced/pizza-cli/issues/22)) ([ad187a9](https://github.com/open-sauced/pizza-cli/commit/ad187a9f3229e41785a09130132a799378c04528))
* Http Client for accessing OpenSauced API client ([#23](https://github.com/open-sauced/pizza-cli/issues/23)) ([ec2b357](https://github.com/open-sauced/pizza-cli/commit/ec2b35789a2864d38bf63e0ec1a3b68393a34e9b))
* Leverage the GITHUB_APP_TOKEN for releases ([#32](https://github.com/open-sauced/pizza-cli/issues/32)) ([e0a25e0](https://github.com/open-sauced/pizza-cli/commit/e0a25e003e89a7a5173ecaae12366922365243c9))
* npm i -g pizza ([73291d1](https://github.com/open-sauced/pizza-cli/commit/73291d13d632b709f2583d834aefe6ad758de8d7))
* repo-query support ([199cfd7](https://github.com/open-sauced/pizza-cli/commit/199cfd7b04e1e1683cce5abc08c57bbef01644f6))
* update bin name release.yaml ([6b21cb8](https://github.com/open-sauced/pizza-cli/commit/6b21cb84f88f75467ce6f270e136dfca5e462d23))


### 🔁 Continuous Integration

* Update @open-sauced/release@2.2.1 and compliance.yaml ([#33](https://github.com/open-sauced/pizza-cli/issues/33)) ([146b6b7](https://github.com/open-sauced/pizza-cli/commit/146b6b7485a0f33090a4ccefd23624f9aa0df085))
