# tlm
Tiny login manager.

This is a login manager designed to be as small (in lines of code) as possible, for flexbility and extensibility. Uses [charm libraries](https://charm.sh/) for visuals. I'm not currently sure how to implement it as an *actual* login manager (à la GDM, XDM) but currently it can be used as e.g. a lock screen.

## How it works
tlm uses PAM with the `login` context for user authentication. Therefore it should work on any device that uses PAM and has a PAM service configuration for `login`.
