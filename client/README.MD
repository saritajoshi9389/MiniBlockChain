https://osxfuse.github.io/
https://github.com/osxfuse/osxfuse/wiki/Mount-options
sudo sysctl macfuse.tunables.admin_group # get
macfuse.tunables.admin_group: 80
$ sudo sysctl -w macfuse.tunables.admin_group=81 # set
macfuse.tunables.admin_group: 80 -> 81