# Changelog

## v0.2.0 (April 14, 2025)

### New Features

* Introduce `fmt` command (7223b43)
* Introduce `check` command (526f026)
* Allow expressions on the `run` command (02758e1)
* List devices (a7af89d)
* Various blocks, corresponding expressions added:
    * `Title` (58c3783)
    * `Clear` (94a151a)
    * `Blur` (3bee9a2)
    * `Stop` (a722815)
    * `Reload` (cd9c65a)
    * `Sleep` (36fe3f0)
    * `NavigateBack` (828d856)
    * `NavigateForward` (14c4661)
    * `Screenshot` (1eb3386)
        * Support scale with screenshot (feef2b1)
    * `DoubleClick` (98b852e)
    * `Flow` (88e06cb)
    * `Navigate` (a3ff68e)
    * `Click` and `WaitVisible` (fdd4d0f)
        * Support selector block within click, wait (b1e2bed)
* Device emulation (2aa7870)

### Bug fixes and improvements

* Report extended version info (b9f5de0)
* Allow additional query options on click, etc. (3991f37)
* Validate identifier names (2e71285)
* Incremental refactorings:
    * Improve late binding of tasks specified from expressions (a3ff68e)
    * Simplify parsing of selector enums (d3a8377)
    * More reducers in selectors handling (b62ae07)
    * Further simplification via reducers (1b5e42c)
    * Use task reducers on recently added task decoders (3161af9)
    * Rename selector{Action,Task} (e7c612b)
    * Introduce automation pkg and reorganize (a840bd0)
    * Encapsulate output variable as higher order function (c16b03e)
    * Enapsulate automation result as context value (86f9c95)
    * Simplify withBinding to benefit from upstream arg checks (d6557ac)
    * Introduce contextual pattern; relocate setup func (b224187)
    * Parser tests (fdd4d0f)
    * Make URL in navigate expr; propagate expressions (b28fe10)
    * Split config selector blocks into own file (6f8f304)
    * Rename Allocator.URL to BrowserURL (848b453)
    * Bind: Specify chrome as a binder (3f36a2e)
    * Use bind extension in order to simplify evaluator binding (e39c6b7)
    * Introduce Engine enumeration, flag; Binder rework into interface (edbddbd)
* Chores:
    * Configuration and dependency updates:
        * Various configuration and dependency updates (1fb8a92)
        * Update dependent versions (c8ea31a)
        * Bump goreleaser action (95051dd)
            * Update GoReleaser configuration (75f71e0)
        * Bump actions/setup-go-5 (4b7719c)
        * Bump goreleaser/goreleaser-action from 3 to 6 (c325046)
        * Bump actions/setup-go from 3 to 5 (f49307a)
        * Update dependent versions (0053da7)
        * Bump Go version in CI (6a5684c)
        * Bump actions/checkout-4 (4016ee1)
        * Update dependent versions (c9a7a00)
        * Bump actions/checkout from 3 to 4 (d542d70)
        * Add Dependabot configuration (b6dec4f)
        * Update engineering platform (1670a4a)
        * Update dependent versions; go1.20 (4c22ee2)
        * GitHub CI configuration (bd91b00)
        * Update engineering platform (7e7ca73)
        * Update dependent versions (8969d19)
        * Update to use modern go tools configuration; lint task in Makefile (f62f41a)
    * Bug and style fixes: Found via linter (e74440a)
    * Apply Go modernizations (649839c)
    * Format Autogun files (aeb5d4b)
    * Update radstub (4e1e97d)


## v0.1.0 (July 23, 2022)

* Initial version :sunrise:
