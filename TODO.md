* Figure out a good way to cast between the javascript object representation of a proto and the actual prototype which does not include properties
* Refactor the service module into its own app modules which actually creates the proper app object and services
* Fix typescript error referring to ref types
* A Job can be deleted without removing it from the formula (Test all updates and removals)
* Prevent create API token from overrolling
* Add a test to make sure removing a job added during a formula update appears both in the job and in the formula and removing works as intended
