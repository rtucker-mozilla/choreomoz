# Application Update Script Architecture
### Script Locations
TBD  

## Update Script Logic

* Standard Unix Codes
    * 0
        * Script Ran Successfully 100%
    * 1-254
        * Some level of imperfection to decide if further scripts should be ran
    * 255
        * Script failed completely

### Update Script Types
* pre-update
* update
* post-update

## Update Script Parameters
* Script-Type
    * 3 Possible Values
        * pre-update
        * update
        * post-update
    * pre-update
        * These will be the first to run and all must complete successfully before the update level scripts fun
    * update
        * These will be the actual update scripts.
    * post-update
        * These will be ran after the update and can potentially involve a reboot followed by futher logic
* Script-Order
    * THIS HAS CURRENTLY NOT BEEN IMPLEMENTED
    * SCRIPT ORDER BASED ON FILENAME
    * An integer value of when the script should run, starting at 0
    * This can be ommitted and Script-Required will be used
    * If both are omitted, then the updater will just run the script at an undetermined time.
* Script-Required
    * THIS HAS CURRENTLY NOT BEEN IMPLEMENTED
    * The name of a script that must have completed successfully before this script is ran
* Script-Exit-Code-Reboot
    * If this is set and the script returns the exit code as provided, the system will
    save it's state to a local file and reboot
* Script-Exit-Max
    * THIS HAS CURRENTLY NOT BEEN IMPLEMENTED
    * The maximum error code of the previous script which will allow this to run.
    * If Script-Exit-Max is set to 128 and the previous script had an exit code of 128 or less then this script will run
* Script-Description
    * A description of the script that is running for logging purposes
* Script-Docs
    * A http link to the documentation for this script

### Update Script Headers Example
    ### BEGIN AUTO UPDATE INFO
    # Script-Type:        pre-update
    # Script-Order:       1
    # Script-Required:    remove_from_loadbalancer.sh
    # Script-Exit-Max:    128
    # Script-Description: Description of what this script does
    # Script-Docs:        https://domain.com/update_script_info.html
    ### END AUTO UPDATE INFO