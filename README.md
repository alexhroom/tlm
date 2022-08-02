# tlm
Tiny login manager.

This is a login manager designed to be as small as possible. I don't like many of the existing options (XDM and GDM for example are far too big and bulky) so this one aims to be as small as possible. Uses [charm libraries](https://charm.sh/) for visuals.

## How it works
tlm uses PAM with the `login` context for user authentication. Therefore it should work on any device that uses PAM and has a PAM service configuration for `login`.
