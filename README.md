# PicoInit
Init system / service manager for seedboxes and non-root user accounts.

**NOTE: Don't rely on this for production software, this was written in two hours with ZERO research on how init systems work.** 

## Need
Tl;dr: Seedbox didnt have systemd or any init system accessible to users and relied on weird cron hack to launch applications in GNU Screen, and this hack called PicoInit is better than said cron hack.

## Usage
* Compile the executable or download from the releases here.
* Make a new folder and put the executable in it.
* Create a file called `picoinit_config.json` and a folder called `logs` in the same folder as your executable.
* Fillup `picoinit_config.json`
* Find a way to run this in the background. Use screen/tmux if you cant find anything.
* Your service logs will be stored in the logs folder.
* There is no way to add or manually restart a service without killing and restarting PicoInit.

Example **`picoinit_config.json`** provided in this repo.

Key | Value 
--- | --- 
`name` | A name for your service without spaces.
`workdir` | Working directory (absolute path) from where the service should be run. Leave it empty if not required.
`cmd` | Your service's command to execute. It shoud either be an absolute or the bin should be in your path. 
`restart` | Restart policy for your service. One of `always`, `never` and `on_error`.
